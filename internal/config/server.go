// Содержит переменные окружение и флаги для описания конфигурации сервера сбора метрик.
package config

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"

	"github.com/caarlos0/env"
)

type ServerConf struct {
	ServerAddr      string `env:"ADDRESS" json:"address"`
	LogLevel        string `env:"LOG_LEVEL"`
	StoreInterval   int64  `env:"STORE_INTERVAL" json:"store_interval"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"store_file"`
	Restore         bool   `env:"RESTORE" json:"restore"`
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"`
	HashKey         string `env:"KEY"`
	IsActivePprof   bool   `env:"IS_ACTIVE_PPROF" envDefault:"false"`
	CryptoKey       string `env:"CRYPTO_KEY" json:"crypto_key"`
	ConfigPath      string `env:"CONFIG"`
	TrustedSubnet   string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
}

func NewServerConf() (*ServerConf, error) {
	conf := &ServerConf{}
	flag.StringVar(&conf.ConfigPath, "c", "internal/config/server_config.json", "Path to json config file")
	err := conf.loadConfigFile()
	if err != nil {
		return nil, err
	}

	flag.StringVar(&conf.ServerAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&conf.LogLevel, "l", "info", "log level")
	flag.StringVar(&conf.DatabaseDSN, "d", "", "setup database dsn connection settings")
	flag.Int64Var(&conf.StoreInterval, "i", 300, "store interval")
	flag.StringVar(&conf.FileStoragePath, "f", "/tmp/metrics-db.json", "file storage path")
	flag.BoolVar(&conf.Restore, "r", true, "load exist data at server start")
	flag.StringVar(&conf.HashKey, "k", "test", "hash key for check request")
	flag.BoolVar(&conf.IsActivePprof, "p", false, "enable pprof")
	flag.StringVar(&conf.CryptoKey, "crypto-key", "/Users/ruslan.golovizin/Projects/practicum/keys/private.pem", "path to private key")
	flag.StringVar(&conf.TrustedSubnet, "t", "", "trusted subnet")
	flag.Parse()

	err = env.Parse(conf)
	if err != nil {
		return nil, err
	}

	return conf, err
}

func (conf *ServerConf) loadConfigFile() error {
	if conf.ConfigPath == "" {
		return nil
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	filePath := filepath.Join(currentDir, string(os.PathSeparator), conf.ConfigPath)
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
