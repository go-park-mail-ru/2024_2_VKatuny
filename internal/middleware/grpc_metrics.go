package middleware

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/metrics"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	AuthMicroservice     = 1
	CompressMicroservice = 2
)

func MetricsInterceptor(metrics *metrics.Metrics, logger *logrus.Logger, service int) func(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	grpcHandler grpc.UnaryHandler,
) (interface{}, error) {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		grpcHandler grpc.UnaryHandler,
	) (interface{}, error) {
		logger.Debugf("Received gRPC request: %s", info.FullMethod)
		start := time.Now()
		resp, err := grpcHandler(ctx, req)

		end := time.Since(start)
		var errMsg string
		if err != nil {
			errMsg = err.Error()
		}

		switch service {
		case AuthMicroservice:
			metrics.AuthHits.WithLabelValues(info.FullMethod, errMsg).Inc()
			metrics.AuthTimings.WithLabelValues(info.FullMethod).Observe(end.Seconds())
			logger.Debug("Auth metrics wrote")
		case CompressMicroservice:
			metrics.CompressHits.WithLabelValues(info.FullMethod, errMsg).Inc()
			metrics.CompressTimings.WithLabelValues(info.FullMethod).Observe(end.Seconds())
			logger.Debug("Compress metrics wrote")
		default:
			logger.Errorf("Unknown microservice code %d", service)
		}
		
		return resp, err
	}
}
