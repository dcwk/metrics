package main

import (
	"github.com/dcwk/metrics/internal/client"
	"github.com/dcwk/metrics/internal/config"
	"log"
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
