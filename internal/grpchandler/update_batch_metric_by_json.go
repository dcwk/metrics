package grpchandler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dcwk/metrics/internal/logger"
	"github.com/dcwk/metrics/internal/models"
)

func (s *MetricsServer) UpdateBatchMetricByJSON(
	ctx context.Context,
	req *UpdateBatchMetricByJSONRequest,
) (*UpdateBatchMetricByJSONResponse, error) {
	var metricsData []models.Metrics
	logger.Log.Info(fmt.Sprintf("metrics have been collected for saving: %s", req.Metrics))
	if err := json.Unmarshal([]byte(req.Metrics), &metricsData); err != nil {
		return nil, err
	}

	metricsList := &models.MetricsList{
		List: metricsData,
	}

	if err := s.Storage.AddMetricsAtBatchMode(metricsList); err != nil {
		return nil, err
	}

	return &UpdateBatchMetricByJSONResponse{}, nil
}
