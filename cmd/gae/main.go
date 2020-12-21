package main

import (
	log "github.com/sirupsen/logrus"
	"os"

	"moh.gov.bz/mch/emtct/internal/config"
	"moh.gov.bz/mch/emtct/internal/server"
)

func main() {
	cnf := config.AppConf{
		EmtctDb:   config.DbConf{},
		Auth:      config.AuthConf{},
		AcsisDb:   config.DbConf{},
		ProjectId: os.Getenv("PROJECT_ID"),
	}
	log.WithFields(log.Fields{
		"conf": cnf,
	}).Info("Starting server")
	server.NewServer(cnf)
}
