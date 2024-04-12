package main

import (
	"log"

	"github.com/dcwk/metrics/internal/client"
	"github.com/dcwk/metrics/internal/config"
)

func main() {
	conf, err := config.NewClientConf()
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Run(conf); err != nil {
		log.Fatal(err)

		return
	}
}
