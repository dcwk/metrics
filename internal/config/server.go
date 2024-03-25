package config

import (
	"flag"

	"github.com/caarlos0/env"
)

type ServerConf struct {
	ServerAddr string `env:"ADDRESS"`
	LogLevel   string `env:"LOG_LEVEL"`
}

func NewServerConf() (*ServerConf, error) {
	conf := &ServerConf{}

	flag.StringVar(&conf.ServerAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&conf.LogLevel, "l", "info", "log level")
	flag.Parse()

	err := env.Parse(conf)
	if err != nil {
		return nil, err
	}

	return conf, err
}
