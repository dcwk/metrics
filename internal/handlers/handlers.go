package handlers

import (
	"github.com/dcwk/metrics/internal/config"
	"github.com/dcwk/metrics/internal/storage"
)

type Handlers struct {
	Storage    storage.DataKeeper
	ServerConf *config.ServerConf
	ClientConf *config.ClientConf
}
