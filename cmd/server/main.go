package main

import (
	"github.com/dcwk/metrics/internal/config"
	"github.com/dcwk/metrics/internal/server"
)

func main() {
	conf := config.NewServerConf()
	server.Run(conf)
}
