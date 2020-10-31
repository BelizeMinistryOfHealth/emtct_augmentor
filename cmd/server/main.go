package main

import (
	"flag"
	"fmt"
	"os"

	"moh.gov.bz/mch/emtct/internal/config"
	"moh.gov.bz/mch/emtct/internal/http"
)

func main() {
	var confFile string
	flag.StringVar(&confFile, "c", "", "Specify configuration file.")
	flag.Parse()
	if len(confFile) == 0 {
		fmt.Errorf("please specify the configuration file using the -c flag")
		os.Exit(1)
	}
	cnf, err := config.ReadConf(confFile)
	if err != nil {
		fmt.Errorf("could not parse the configuration file")
	}
	http.NewServer(*cnf)
}
