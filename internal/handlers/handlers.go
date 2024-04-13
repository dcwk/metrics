package handlers

import (
	"database/sql"

	"github.com/dcwk/metrics/internal/storage"
)

type Handlers struct {
	Storage storage.DataKeeper
	DB      *sql.DB
}
