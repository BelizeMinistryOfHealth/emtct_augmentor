package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	"moh.gov.bz/mch/emtct/internal/app"
	"moh.gov.bz/mch/emtct/internal/app/api"
	"moh.gov.bz/mch/emtct/internal/config"
	"moh.gov.bz/mch/emtct/internal/db"
)

func RegisterHandlers(ctx context.Context, cnf config.AppConf) (*mux.Router, error) {
	firestoreDb, err := db.NewFirestore(ctx, cnf.ProjectId)
	if err != nil {
		log.WithError(err).Error("firestore connection failed")
		os.Exit(-1)
	}

	app := app.App{
		ProjectID:       cnf.ProjectId,
		Firestore:       firestoreDb,
		FirestoreApiKey: cnf.FirestoreApiKey,
	}
	router, err := api.API(app)
	if err != nil {
		return nil, fmt.Errorf("failed to register handlers: %w", err)
	}
	log.Infof("Initiated App: %+v", app)

	return router, nil
}

func NewServer(cnf config.AppConf) {
	firestoreContext := context.Background()
	r, err := RegisterHandlers(firestoreContext, cnf)
	if err != nil {
		log.WithFields(log.Fields{
			"appConf": cnf,
		}).WithError(err).Error("failed to register handlers")
		os.Exit(-1)
	}
	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", cnf.Port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 60,
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
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	// Create a deadline to wait for.
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
