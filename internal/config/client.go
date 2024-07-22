// Содержит переменные окружение и флаги для описания конфигурации агента для отправки метрик
package config

import (
	"flag"

	"github.com/caarlos0/env"
)

type ClientConf struct {
	ServerAddr     string `env:"ADDRESS"`
	ReportInterval int64  `env:"REPORT_INTERVAL"`
	PollInterval   int64  `env:"POLL_INTERVAL"`
	LogLevel       string `env:"LOG_LEVEL"`
	HashKey        string `env:"KEY"`
	CryptoKey      string `env:"CRYPTO_KEY"`
}

func NewClientConf() (*ClientConf, error) {
	conf := &ClientConf{}

	flag.StringVar(&conf.ServerAddr, "a", "localhost:8080", "metrics server address")
	flag.Int64Var(&conf.ReportInterval, "r", 10, "sending frequency interval")
	flag.Int64Var(&conf.PollInterval, "p", 2, "metrics reading frequency")
	flag.StringVar(&conf.LogLevel, "l", "info", "log level")
	flag.StringVar(&conf.HashKey, "k", "test", "hash key")
	flag.StringVar(&conf.CryptoKey, "crypto-key", "", "path to public key")
	flag.Parse()

	err := env.Parse(conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
