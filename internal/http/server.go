package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/uris77/auth0"

	"moh.gov.bz/mch/emtct/internal/config"
	"moh.gov.bz/mch/emtct/internal/db"

	_ "github.com/lib/pq"
)

type App struct {
	Db      *db.EmtctDb
	AcsisDb *db.AcsisDb
}

func RegisterHandlers(cnf config.AppConf) *mux.Router {
	jwkUrl := cnf.Auth.JwkUrl
	iss := cnf.Auth.Issuer
	aud := cnf.Auth.Audience

	// Instantiate an aut0 client with a Cache with a key capacity of
	// 60 tokens and a ttl of 24 hours.
	auth0Client := auth0.NewAuth0(60, 518400)

	emtctDb, err := db.NewConnection(&cnf.EmtctDb)
	if err != nil {
		log.WithError(err).Error("failed to open connection to database")
		panic("failed to open connection to database")
	}
	acsisDb, err := db.NewAcsisConnection(&cnf.AcsisDb)
	if err != nil {
		log.WithError(err).Error("failed to open connection to acsis database")
		panic("failed to open connection to acsis database")
	}

	app := App{Db: emtctDb, AcsisDb: acsisDb}

	r := mux.NewRouter()
	authHandlers := NewChain(EnableCors(), VerifyToken(jwkUrl, aud, iss, auth0Client))

	//Infant Syphilis Screening
	r.HandleFunc("/patient/{motherId}/infant/{infantId}/syphilisScreenings",
		authHandlers.Then(app.InfantSyphilisScreeningHandler)).
		Methods(http.MethodGet, http.MethodOptions)

	fs := http.FileServer(http.Dir("./front_end/build/"))
	r.PathPrefix("/").Handler(http.StripPrefix("/", fs))
	staticFs := http.FileServer(http.Dir("./front_end/build/static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFs))

	return r
}

func NewServer(cnf config.AppConf) {
	r := RegisterHandlers(cnf)
	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Println("Starting server on port 8080")
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	wait := time.Duration(30)
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
