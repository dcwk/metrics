package storage

import (
	"errors"
	"strconv"
	"strings"
)

type Storage interface {
	AddGauge(name string, value float64)
	AddCounter(name string, value int64)
	GetGauge(name string) (string, error)
	GetCounter(name string)
}

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewStorage() MemStorage {
	return MemStorage{
		gauge:   map[string]float64{},
		counter: map[string]int64{},
	}
}

func (ms *MemStorage) AddGauge(name string, value string) error {
	convertedVal, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil {
		return errors.New("unsupported gauge value")
	}

	ms.gauge[name] = convertedVal

	return nil
}

func (ms *MemStorage) GetGauge(name string) (float64, error) {
	if ms.gauge[name] == 0 {
		return 0, errors.New("gauge not found")
	}

	return ms.gauge[name], nil
}

func (ms *MemStorage) AddCounter(name string, value string) error {
	convertedVal, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil {
		return errors.New("unsupported counter value")
	}

	ms.counter[name] = convertedVal

	return nil
}

func (ms *MemStorage) GetCounter(name string) (int64, error) {
	if ms.counter[name] == 0 {
		return 0, errors.New("counter not found")
	}

	return ms.counter[name], nil
}
