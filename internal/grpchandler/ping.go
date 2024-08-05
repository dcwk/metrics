package grpchandler

import (
	"context"
	"fmt"

	"github.com/dcwk/metrics/internal/logger"
)

func (s *MetricsServer) pingHandler(ctx context.Context, pingRequest *PingRequest) (*PingResponse, error) {
	logger.Log.Info("Fired ping grpc handler")
	err := s.Storage.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't connect to db: %w", err)
	}

	return &PingResponse{}, nil
}
