package grpchandler

import (
	"github.com/dcwk/metrics/internal/config"
	"github.com/dcwk/metrics/internal/storage"
)

type MetricsServer struct {
	Storage storage.DataKeeper
	Conf    *config.ServerConf
	UnimplementedMetricsServiceServer
}
