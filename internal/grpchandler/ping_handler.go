package grpchandler

import (
	"context"
	"fmt"
)

func (s *MetricsServer) pingHandler(ctx context.Context, pingRequest *PingRequest) (*PingResponse, error) {
	err := s.Storage.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't connect to db: %w", err)
	}

	return &PingResponse{}, nil
}
