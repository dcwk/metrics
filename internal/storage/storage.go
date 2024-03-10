package storage

import (
	"errors"
	"strconv"
	"strings"
	"sync"
)

var stor *MemStorage

type Storage interface {
	AddGauge(name string, value float64)
	AddCounter(name string, value int64)
	GetGauge(name string) (string, error)
	GetCounter(name string)
}

type Gauge struct {
	mx sync.RWMutex
	m  map[string]float64
}

type Counter struct {
	mx sync.RWMutex
	m  map[string]int64
}
type MemStorage struct {
	Gauge
	Counter
}

func NewStorage() *MemStorage {
	if stor == nil {
		stor = &MemStorage{
			Gauge{m: map[string]float64{}},
			Counter{m: map[string]int64{}},
		}
	}

	return stor
}

func (ms *MemStorage) AddGauge(name string, value string) error {
	ms.Gauge.mx.RLock()
	defer ms.Gauge.mx.RUnlock()

	convertedVal, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil {
		return errors.New("unsupported gauge value")
	}

	ms.Gauge.m[name] = convertedVal

	return nil
}

func (ms *MemStorage) GetGauge(name string) (float64, error) {
	if ms.Gauge.m[name] == 0 {
		return 0, errors.New("gauge not found")
	}

	return ms.Gauge.m[name], nil
}

func (ms *MemStorage) GetAllGauges() map[string]float64 {
	return ms.Gauge.m
}

func (ms *MemStorage) AddCounter(name string, value string) error {
	ms.Counter.mx.RLock()
	defer ms.Counter.mx.RUnlock()

	convertedVal, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil {
		return errors.New("unsupported counter value")
	}

	ms.Counter.m[name] = convertedVal

	return nil
}

func (ms *MemStorage) GetCounter(name string) (int64, error) {
	if ms.Counter.m[name] == 0 {
		return 0, errors.New("counter not found")
	}

	return ms.Counter.m[name], nil
}

func (ms *MemStorage) GetAllCounters() map[string]int64 {
	return ms.Counter.m
}
