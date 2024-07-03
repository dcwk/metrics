package main

import (
	"fmt"
	"log"
	_ "net/http/pprof"

	"github.com/dcwk/metrics/internal/config"
	"github.com/dcwk/metrics/internal/server"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {
	conf, err := config.NewServerConf()
	if err != nil {
		log.Fatal(err)
	}

	buildInfo()
	server.Run(conf)
}

func buildInfo() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Println(fmt.Sprintf("Build version: %s", buildVersion))
	fmt.Println(fmt.Sprintf("Build date: %s", buildDate))
	fmt.Println(fmt.Sprintf("Build commit: %s", buildCommit))
}
