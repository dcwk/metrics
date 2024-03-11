package config

import (
	"flag"
	"github.com/caarlos0/env"
	"log"
)

type ServerConf struct {
	ServerAddr string `env:"ADDRESS"`
}

func NewServerConf() ServerConf {
	conf := ServerConf{}

	flag.StringVar(&conf.ServerAddr, "a", ":8080", "address and port to run server")
	flag.Parse()

	err := env.Parse(&conf)
	if err != nil {
		log.Fatal(err)
	}

	return conf
}
