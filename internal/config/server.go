package config

import (
	"flag"

	"github.com/caarlos0/env"
)

type ServerConf struct {
	ServerAddr      string `env:"ADDRESS"`
	LogLevel        string `env:"LOG_LEVEL"`
	StoreInterval   int64  `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

func NewServerConf() (*ServerConf, error) {
	conf := &ServerConf{}

	flag.StringVar(&conf.ServerAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&conf.LogLevel, "l", "info", "log level")
	flag.StringVar(&conf.DatabaseDSN, "d", "host=localhost user=test password=videos dbname=videos sslmode=disable", "setup database dsn connection settings")
	flag.Int64Var(&conf.StoreInterval, "i", 300, "store interval")
	flag.StringVar(&conf.FileStoragePath, "f", "/tmp/metrics-db.json", "file storage path")
	flag.BoolVar(&conf.Restore, "r", true, "load exist data at server start")
	flag.Parse()

	err := env.Parse(conf)
	if err != nil {
		return nil, err
	}

	return conf, err
}
