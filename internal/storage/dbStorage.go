package storage

import (
	"context"
	"database/sql"
	"os"
	"sync"

	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/models"
	"github.com/pressly/goose"
)

type DatabaseStorage struct {
	DB *sql.DB
	mu sync.Mutex
}

func NewDBStorage(db *sql.DB) (*DatabaseStorage, error) {
	dbs := &DatabaseStorage{
		DB: db,
		mu: sync.Mutex{},
	}

	dbs.mu.Lock()
	defer dbs.mu.Unlock()

	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	pwd = pwd + "/../../migrations"
	if err := goose.Up(db, pwd); err != nil {
		logger.Log.Error("Can't apply migrations")
		return nil, err
	}

	return dbs, nil
}

func (dbs *DatabaseStorage) AddGauge(name string, value float64) error {
	dbs.mu.Lock()
	defer dbs.mu.Unlock()

	_, err := dbs.DB.Exec(
		"INSERT INTO gauges (id, value) VALUES ($1, $2) ON CONFLICT(id) DO UPDATE SET value=$2",
		name,
		value,
	)
	if err != nil {
		return err
	}

	return nil
}

func (dbs *DatabaseStorage) GetGauge(name string) (float64, error) {
	dbs.mu.Lock()
	defer dbs.mu.Unlock()
	var gaugeValue float64

	row := dbs.DB.QueryRow("SELECT value FROM gauges WHERE id=$1", name)
	if err := row.Scan(&gaugeValue); err != nil {
		return gaugeValue, err
	}

	return gaugeValue, nil
}

func (dbs *DatabaseStorage) GetAllGauges() map[string]float64 {
	gaugesMap := make(map[string]float64, 0)
	dbs.mu.Lock()
	defer dbs.mu.Unlock()

	rows, err := dbs.DB.Query("SELECT id, value FROM gauges")
	if err != nil {
		return gaugesMap
	}
	defer rows.Close()

	for rows.Next() {
		var metrics models.Metrics
		err := rows.Scan(&metrics.ID, &metrics.Value)
		if err != nil {
			return gaugesMap
		}

		gaugesMap[metrics.ID] = *metrics.Value
	}

	if err := rows.Err(); err != nil {
		return gaugesMap
	}

	return gaugesMap
}

func (dbs *DatabaseStorage) AddCounter(name string, value int64) error {
	dbs.mu.Lock()
	defer dbs.mu.Unlock()

	_, err := dbs.DB.Exec(
		"INSERT INTO counters AS c (id, delta) VALUES ($1, $2) ON CONFLICT(id) DO UPDATE SET delta=c.delta + $2",
		name,
		value,
	)
	if err != nil {
		return err
	}

	return nil
}

func (dbs *DatabaseStorage) GetCounter(name string) (int64, error) {
	dbs.mu.Lock()
	defer dbs.mu.Unlock()
	var delta int64

	row := dbs.DB.QueryRow("SELECT delta FROM counters WHERE id=$1", name)
	if err := row.Scan(&delta); err != nil {
		return delta, err
	}

	return delta, nil
}

func (dbs *DatabaseStorage) GetAllCounters() map[string]int64 {
	countersMap := make(map[string]int64, 0)
	dbs.mu.Lock()
	defer dbs.mu.Unlock()

	rows, err := dbs.DB.Query("SELECT id, value FROM gauges")
	if err != nil {
		return countersMap
	}
	defer rows.Close()

	for rows.Next() {
		var metrics models.Metrics
		err := rows.Scan(&metrics.ID, &metrics.Delta)
		if err != nil {
			return countersMap
		}

		countersMap[metrics.ID] = *metrics.Delta
	}

	if err := rows.Err(); err != nil {
		return countersMap
	}

	return countersMap
}

func (dbs *DatabaseStorage) AddMetricsAtBatchMode(metricsList *models.MetricsList) error {
	tx, err := dbs.DB.Begin()
	if err != nil {
		return err
	}

	for _, v := range metricsList.List {
		if v.MType == models.Gauge {
			if err := dbs.AddGauge(v.ID, *v.Value); err != nil {
				return tx.Rollback()
			}
		}

		if v.MType == models.Counter {
			if err := dbs.AddCounter(v.ID, *v.Delta); err != nil {
				return tx.Rollback()
			}
		}
	}

	return tx.Commit()
}

func (dbs *DatabaseStorage) Ping(ctx context.Context) error {
	err := dbs.DB.PingContext(ctx)
	if err != nil {
		return err
	}

	return nil
}
