package service

import (
	"flag"
	"os"

	"github.com/realbucksavage/robin/pkg/log"
)

func Run() {
	log.L.Infof("Program started with command arguments: %v", os.Args)

	cfg := flag.String("config", "./robinconf.yaml", "Configuration file to read from")
	loggingLevel := flag.String("logging-level", "INFO", "")
	flag.Parse()

	log.SetLevel(*loggingLevel)

	config, err := readConfig(*cfg)
	if err != nil {
		log.L.Fatal(err)
	}

	cmdService(config)
}
