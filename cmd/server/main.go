package main

import (
	"log"
	_ "net/http/pprof"

	"github.com/dcwk/metrics/internal/config"
	"github.com/dcwk/metrics/internal/server"
	"github.com/dcwk/metrics/internal/utils/info"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {
	conf, err := config.NewServerConf()
	if err != nil {
		log.Fatal(err)
	}

	info.BuildInfo(buildVersion, buildDate, buildCommit)
	server.Run(conf)
}
