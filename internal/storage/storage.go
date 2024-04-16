package storage

import (
	"errors"
	"sync"

	"github.com/dcwk/metrics/internal/models"
	"github.com/mailru/easyjson"
)

type DataKeeper interface {
	AddGauge(name string, value float64) error
	GetGauge(name string, allowZeroVal bool) (float64, error)
	GetAllGauges() map[string]float64
	AddCounter(name string, value int64) error
	GetCounter(name string, allowZeroVal bool) (int64, error)
	GetAllCounters() map[string]int64
	GetJSONMetrics() (string, error)
	SaveMetricsList(metricsList *models.MetricsList)
}

type Gauge struct {
	gaugeMx sync.Mutex
}

type Counter struct {
	counterMx sync.Mutex
}

type MemStorage struct {
	mu      sync.Mutex
	gauge   map[string]float64
	counter map[string]int64
}

func NewStorage() *MemStorage {
	return &MemStorage{
		mu:      sync.Mutex{},
		gauge:   make(map[string]float64, 1000),
		counter: make(map[string]int64, 1000),
	}
}

func (ms *MemStorage) AddGauge(name string, value float64) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.gauge[name] = value

	return nil
}

func (ms *MemStorage) GetGauge(name string, allowZeroVal bool) (float64, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if ms.gauge[name] == 0 && !allowZeroVal {
		return 0, errors.New("gauge not found")
	}

	val, ok := ms.gauge[name]
	if !ok {
		return 0, errors.New("failed to get metric")
	}

	return val, nil
}

func (ms *MemStorage) GetAllGauges() map[string]float64 {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	return ms.gauge
}

func (ms *MemStorage) AddCounter(name string, value int64) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.counter[name] += value

	return nil
}

func (ms *MemStorage) GetCounter(name string, allowZeroVal bool) (int64, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if ms.counter[name] == 0 && !allowZeroVal {
		return 0, errors.New("counter not found")
	}

	val, ok := ms.counter[name]
	if !ok {
		return 0, errors.New("failed to get metric")
	}

	return val, nil
}

func (ms *MemStorage) GetAllCounters() map[string]int64 {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	return ms.counter
}

func (ms *MemStorage) GetJSONMetrics() (string, error) {
	metricsList := models.MetricsList{}
	for k, v := range ms.gauge {
		metric := models.Metrics{
			ID:    k,
			MType: models.Gauge,
			Value: &v,
		}
		metricsList.List = append(metricsList.List, metric)
	}

	for k, v := range ms.counter {
		metric := models.Metrics{
			ID:    k,
			MType: models.Counter,
			Delta: &v,
		}
		metricsList.List = append(metricsList.List, metric)
	}

	metricsListJSON, err := easyjson.Marshal(&metricsList)
	if err != nil {
		return "", err
	}

	return string(metricsListJSON), nil
}

func (ms *MemStorage) SaveMetricsList(metricsList *models.MetricsList) {
	for _, v := range metricsList.List {
		if v.MType == models.Gauge {
			_ = ms.AddGauge(v.ID, *v.Value)
		}

		if v.MType == models.Counter {
			_ = ms.AddCounter(v.ID, *v.Delta)
		}
	}
}
