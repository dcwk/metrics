package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dcwk/metrics/internal/config"
	"github.com/dcwk/metrics/internal/server"
)

func main() {
	conf, err := config.NewServerConf()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop
		cancel()
	}()

	server.Run(ctx, conf)
}
