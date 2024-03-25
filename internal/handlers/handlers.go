package handlers

import "github.com/dcwk/metrics/internal/storage"

const (
	gauge   = "gauge"
	counter = "counter"
)

type Handlers struct {
	Storage storage.DataKeeper
}
