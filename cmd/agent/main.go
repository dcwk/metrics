package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dcwk/metrics/internal/client"
	"github.com/dcwk/metrics/internal/config"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {
	conf, err := config.NewClientConf()
	if err != nil {
		log.Fatal(err)
	}
	buildInfo()

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
