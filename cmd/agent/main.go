package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dcwk/metrics/internal/client"
	"github.com/dcwk/metrics/internal/config"
	"github.com/dcwk/metrics/internal/utils"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {
	conf, err := config.NewClientConf()
	if err != nil {
		log.Fatal(err)
	}

	utils.BuildInfo(buildVersion, buildDate, buildCommit)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop

		cancel()
	}()

	if err := client.Run(ctx, conf); err != nil {
		log.Fatal(err)

		return
	}
}
