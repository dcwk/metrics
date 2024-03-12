package storage

import (
	"errors"
	"strconv"
	"strings"
	"sync"
)

type Storage interface {
	AddGauge(name string, value string) error
	GetGauge(name string) (float64, error)
	GetAllGauges() map[string]float64
	AddCounter(name string, value string) error
	GetCounter(name string) (int64, error)
	GetAllCounters() map[string]int64
}

type Gauge struct {
	gaugeMx sync.RWMutex
	gauge   map[string]float64
}

type Counter struct {
	counterMx sync.RWMutex
	counter   map[string]int64
}

type MemStorage struct {
	Gauge
	Counter
}

func NewStorage() *MemStorage {
	return &MemStorage{
		Gauge{gauge: make(map[string]float64, 1000)},
		Counter{counter: make(map[string]int64, 1000)},
	}
}

func (ms *MemStorage) AddGauge(name string, value string) error {
	ms.gaugeMx.Lock()
	defer ms.gaugeMx.Unlock()

	convertedVal, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil {
		return errors.New("unsupported gauge value")
	}

	ms.gauge[name] = convertedVal

	return nil
}

func (ms *MemStorage) GetGauge(name string) (float64, error) {
	ms.gaugeMx.RLock()
	defer ms.gaugeMx.RUnlock()

	if ms.gauge[name] == 0 {
		return 0, errors.New("gauge not found")
	}

	return ms.gauge[name], nil
}

func (ms *MemStorage) GetAllGauges() map[string]float64 {
	ms.gaugeMx.RLock()
	defer ms.gaugeMx.RUnlock()

	return ms.gauge
}

func (ms *MemStorage) AddCounter(name string, value string) error {
	ms.counterMx.Lock()
	defer ms.counterMx.Unlock()

	convertedVal, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil {
		return errors.New("unsupported counter value")
	}

	ms.counter[name] += convertedVal

	return nil
}

func (ms *MemStorage) GetCounter(name string) (int64, error) {
	ms.counterMx.RLock()
	defer ms.counterMx.RUnlock()

	if ms.counter[name] == 0 {
		return 0, errors.New("counter not found")
	}

	return ms.counter[name], nil
}

func (ms *MemStorage) GetAllCounters() map[string]int64 {
	ms.counterMx.RLock()
	defer ms.counterMx.RUnlock()

	return ms.counter
}
