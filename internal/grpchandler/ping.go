package grpchandler

import (
	"context"
	"fmt"

	"github.com/dcwk/metrics/internal/logger"
)

func (s *MetricsServer) PingHandler(ctx context.Context, pingRequest *PingRequest) (*PingResponse, error) {
	logger.Log.Info("the grpc ping request handler is running")
	err := s.Storage.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't connect to db: %w", err)
	}

	return &PingResponse{}, nil
}
