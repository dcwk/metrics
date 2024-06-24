// Пакет service содержит сервисы для упрощения работы с логикой приложения
package service

import (
	"errors"

	"github.com/dcwk/metrics/internal/models"
	"github.com/dcwk/metrics/internal/storage"
)

type MetricsService struct {
	storage storage.DataKeeper
}

func NewMetricsService(storage storage.DataKeeper) *MetricsService {
	return &MetricsService{
		storage: storage,
	}
}

// UpdateMetrics - обновляет метрики в выбранном хранилище
func (ms *MetricsService) UpdateMetrics(metrics *models.Metrics) error {
	switch metrics.MType {
	default:
		return errors.New("type doesn't support")
	case models.Gauge:
		if err := ms.storage.AddGauge(metrics.ID, *metrics.Value); err != nil {
			return err
		}
	case models.Counter:
		if err := ms.storage.AddCounter(metrics.ID, *metrics.Delta); err != nil {
			return err
		}
	}

	return nil
}

// GetMetrics - загружает метрики из выбранного хранилища
func (ms *MetricsService) GetMetrics(metrics *models.Metrics) (*models.Metrics, error) {
	switch metrics.MType {
	default:
		return metrics, errors.New("unsupported type")
	case models.Gauge:
		metricValue, err := ms.storage.GetGauge(metrics.ID)
		if err != nil {
			return metrics, err
		}

		metrics.Value = &metricValue
	case models.Counter:
		metricValue, err := ms.storage.GetCounter(metrics.ID)
		if err != nil {
			return metrics, err
		}
		metrics.Delta = &metricValue
	}

	return metrics, nil
}
