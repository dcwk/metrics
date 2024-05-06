package storage

import (
	"context"
	"database/sql"
	"sync"
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

	_, err := dbs.DB.Exec("CREATE TABLE IF NOT EXISTS public.gauges (id varchar NULL,value double precision NULL)")
	if err != nil {
		return nil, err
	}
	_, err = dbs.DB.Exec("CREATE TABLE IF NOT EXISTS public.counters (id varchar NULL,delta int NULL)")
	if err != nil {
		return nil, err
	}

	return dbs, nil
}

func (dbs *DatabaseStorage) AddGauge(name string, value float64) error {
	dbs.mu.Lock()
	defer dbs.mu.Unlock()

	_, err := dbs.DB.Exec("INSERT INTO gauges (id, value) VALUES ($1, $2)", name, value)
	if err != nil {
		return err
	}

	return nil
}

func (dbs *DatabaseStorage) GetGauge(name string, allowZeroVal bool) (float64, error) {
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
	gauges := make(map[string]float64)
	return gauges
}

func (dbs *DatabaseStorage) AddCounter(name string, value int64) error {
	dbs.mu.Lock()
	defer dbs.mu.Unlock()

	_, err := dbs.DB.Exec("INSERT INTO counters (id, delta) VALUES ($1, $2)", name, value)
	if err != nil {
		return err
	}

	return nil
}

func (dbs *DatabaseStorage) GetCounter(name string, allowZeroVal bool) (int64, error) {
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
	counters := make(map[string]int64)
	return counters
}

func (dbs *DatabaseStorage) Ping(ctx context.Context) error {
	err := dbs.DB.PingContext(ctx)
	if err != nil {
		return err
	}

	return nil
}
