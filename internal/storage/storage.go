package storage

import (
	"errors"
	"sync"

	"github.com/dcwk/metrics/internal/models"
	"github.com/mailru/easyjson"
)

type DataKeeper interface {
	AddGauge(name string, value *float64) error
	GetGauge(name string, allowZeroVal bool) (float64, error)
	GetAllGauges() map[string]float64
	AddCounter(name string, value *int64) error
	GetCounter(name string, allowZeroVal bool) (int64, error)
	GetAllCounters() map[string]int64
	GetJSONMetrics() (string, error)
	SaveMetricsList(metricsList *models.MetricsList)
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

func (ms *MemStorage) AddGauge(name string, value *float64) error {
	ms.gaugeMx.Lock()
	defer ms.gaugeMx.Unlock()

	ms.gauge[name] = *value

	return nil
}

func (ms *MemStorage) GetGauge(name string, allowZeroVal bool) (float64, error) {
	ms.gaugeMx.RLock()
	defer ms.gaugeMx.RUnlock()

	if ms.gauge[name] == 0 && allowZeroVal == false {
		return 0, errors.New("gauge not found")
	}

	return ms.gauge[name], nil
}

func (ms *MemStorage) GetAllGauges() map[string]float64 {
	ms.gaugeMx.RLock()
	defer ms.gaugeMx.RUnlock()

	return ms.gauge
}

func (ms *MemStorage) AddCounter(name string, value *int64) error {
	ms.counterMx.Lock()
	defer ms.counterMx.Unlock()

	ms.counter[name] += *value

	return nil
}

func (ms *MemStorage) GetCounter(name string, allowZeroVal bool) (int64, error) {
	ms.counterMx.RLock()
	defer ms.counterMx.RUnlock()

	if ms.counter[name] == 0 && allowZeroVal == false {
		return 0, errors.New("counter not found")
	}

	return ms.counter[name], nil
}

func (ms *MemStorage) GetAllCounters() map[string]int64 {
	ms.counterMx.RLock()
	defer ms.counterMx.RUnlock()

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
			ms.gauge[v.ID] = *v.Value
		}

		if v.MType == models.Counter {
			ms.counter[v.ID] = *v.Delta
		}
	}
}
