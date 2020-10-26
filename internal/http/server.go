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
	Db *db.EmtctDb
}

func RegisterHandlers() *mux.Router {
	jwkUrl := os.Getenv("EMTCT_JWK_URL")
	iss := os.Getenv("EMTCT_AUTH_ISSUER")
	aud := os.Getenv("EMTCT_AUTH_AUDIENCE")

	// Instantiate an aut0 client with a Cache with a key capacity of
	// 60 tokens and a ttl of 24 hours.
	auth0Client := auth0.NewAuth0(60, 518400)

	//Todo: Read the database configuration from a file
	cnf := config.DbConf{
		DbUsername: "postgres",
		DbPassword: "password",
		DbDatabase: "emtct",
		DbHost:     "localhost",
	}
	db, err := db.NewConnection(&cnf)
	if err != nil {
		log.WithError(err).Error("failed to open connection to database")
		panic("failed to open connection to database")
	}

	app := App{Db: db}

	r := mux.NewRouter()
	authHandlers := NewChain(EnableCors(), VerifyToken(jwkUrl, aud, iss, auth0Client))

	r.HandleFunc("/health", NewChain(EnableCors()).Then(HealthCheck)).Methods(http.MethodGet)
	r.HandleFunc("/test", authHandlers.Then(TestAuth)).Methods(http.MethodGet)
	r.HandleFunc("/patient/{id}", authHandlers.Then(app.RetrievePatient)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/patient/{id}/currentPregnancy", authHandlers.Then(app.FindCurrentPregnancy)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/patient/{id}/currentPregnancy/labResults", authHandlers.Then(app.FindPregnancyLabResults)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/patient/homeVisits", authHandlers.Then(app.PostHomeVisit)).Methods(http.MethodOptions, http.MethodPost)
	r.HandleFunc("/patient/{id}/homeVisits", authHandlers.Then(app.FindHomeVisitsByPatient)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/patient/{patientId}/hivScreenings", authHandlers.Then(app.HivScreeningsByPatientIdHandler)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/patient/homeVisit/{homeVisitId}", authHandlers.Then(app.HomeVisitApi)).Methods(http.MethodGet, http.MethodOptions, http.MethodPut)
	r.HandleFunc("/patient/hivScreening", authHandlers.Then(app.CreateHivScreeningHandler)).Methods(http.MethodOptions, http.MethodPost)

	return r
}

func NewServer() {
	r := RegisterHandlers()
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
