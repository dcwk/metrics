package storage

import "errors"

type Storage interface {
	AddMetric(name string, value string)
	GetMetric(name string) (string, error)
}

type MemStorage struct {
	data map[string]string
}

func NewStorage() MemStorage {
	return MemStorage{
		data: map[string]string{},
	}
}

func (ms *MemStorage) AddMetric(name string, value string) {
	ms.data[name] = value
}

func (ms *MemStorage) GetMetric(name string) (string, error) {
	if ms.data[name] == "" {
		return "", errors.New("metric not found")
	}

	return ms.data[name], nil
}
