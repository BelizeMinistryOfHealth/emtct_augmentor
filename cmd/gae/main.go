package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"

	"moh.gov.bz/mch/emtct/internal/config"
	"moh.gov.bz/mch/emtct/internal/server"
)

func main() {
	cnf := config.AppConf{
		ProjectId:       os.Getenv("PROJECT_ID"),
		FirestoreApiKey: os.Getenv("FIREBASE_API_KEY"),
	}
	log.WithFields(log.Fields{
		"conf": cnf,
	}).Info("Starting server")
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	cnf.Port = port
	server.NewServer(cnf)
}
