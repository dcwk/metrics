// Содержит переменные окружение и флаги для описания конфигурации агента для отправки метрик
package config

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"

	"github.com/caarlos0/env"
)

type ClientConf struct {
	ServerAddr     string `env:"ADDRESS" json:"address"`
	GRPCServerAddr string `env:"GRPC_SERVER_ADDR" json:"grpc_server_addr"`
	ReportInterval int64  `env:"REPORT_INTERVAL" json:"report_interval"`
	PollInterval   int64  `env:"POLL_INTERVAL" json:"poll_interval"`
	LogLevel       string `env:"LOG_LEVEL"`
	HashKey        string `env:"KEY"`
	CryptoKey      string `env:"CRYPTO_KEY" json:"crypto_key"`
	ConfigPath     string `env:"CONFIG_PATH"`
}

func NewClientConf() (*ClientConf, error) {
	conf := &ClientConf{}
	flag.StringVar(&conf.ConfigPath, "c", "../internal/config/client_config.json", "path to json config file")
	err := conf.loadConfigFile()
	if err != nil {
		return nil, err
	}

	flag.StringVar(&conf.ServerAddr, "a", "localhost:8080", "metrics server address")
	flag.StringVar(&conf.GRPCServerAddr, "gs", "localhost:3200", "grpc metrics server address")
	flag.Int64Var(&conf.ReportInterval, "r", 10, "sending frequency interval")
	flag.Int64Var(&conf.PollInterval, "p", 2, "metrics reading frequency")
	flag.StringVar(&conf.LogLevel, "l", "info", "log level")
	flag.StringVar(&conf.HashKey, "k", "test", "hash key")
	flag.StringVar(&conf.CryptoKey, "crypto-key", "/Users/ruslan.golovizin/Projects/practicum/keys/public.pem", "path to public key")
	flag.Parse()

	err = env.Parse(conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func (conf *ClientConf) loadConfigFile() error {
	if conf.ConfigPath == "" {
		return nil
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	filePath := filepath.Join(filepath.Dir(currentDir), string(os.PathSeparator), conf.ConfigPath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, conf)
	if err != nil {
		return err
	}

	return nil
}
