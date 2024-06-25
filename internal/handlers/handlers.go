// Основные use case сервиса сбора метрик
package handlers

import (
	"github.com/dcwk/metrics/internal/storage"
)

type Handlers struct {
	Storage storage.DataKeeper
}
