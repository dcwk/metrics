package main

import (
	"github.com/dcwk/metrics/internal/config"
	"github.com/dcwk/metrics/internal/server"
	"log"
)

func main() {
	conf, err := config.NewServerConf()
	if err != nil {
		log.Fatal(err)
	}

	server.Run(conf)
}
