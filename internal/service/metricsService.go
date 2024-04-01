package service

import (
	"errors"

	"github.com/dcwk/metrics/internal/models"
	"github.com/dcwk/metrics/internal/storage"
)

const (
	gauge   = "gauge"
	counter = "counter"
)

type MetricsService struct {
	storage storage.DataKeeper
}

func NewMetricsService(storage storage.DataKeeper) *MetricsService {
	return &MetricsService{
		storage: storage,
	}
}

func (ms *MetricsService) UpdateMetrics(metrics *models.Metrics) error {
	switch metrics.MType {
	default:
		return errors.New("type doesn't support")
	case gauge:
		if err := ms.storage.AddGauge(metrics.ID, metrics.Value); err != nil {
			return err
		}
	case counter:
		if err := ms.storage.AddCounter(metrics.ID, metrics.Delta); err != nil {
			return err
		}
	}

	return nil
}

func (ms *MetricsService) GetMetrics(metrics *models.Metrics) (*models.Metrics, error) {
	switch metrics.MType {
	default:
		return metrics, errors.New("unsupported type")
	case gauge:
		metricValue, err := ms.storage.GetGauge(metrics.ID)
		if err != nil {
			return metrics, err
		}

		metrics.Value = &metricValue
	case counter:
		metricValue, err := ms.storage.GetCounter(metrics.ID)
		if err != nil {
			return metrics, err
		}
		metrics.Delta = &metricValue
	}

	return metrics, nil
}
