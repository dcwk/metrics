package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dcwk/metrics/internal/client"
	"github.com/dcwk/metrics/internal/config"
)

func main() {
	conf, err := config.NewClientConf()
	if err != nil {
		log.Fatal(err)
	}
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
