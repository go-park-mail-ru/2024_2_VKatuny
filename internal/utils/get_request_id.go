package utils

import (
	"context"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/sirupsen/logrus"
)

func SetLoggerRequestID(ctx context.Context, logger *logrus.Entry) *logrus.Entry {
	requestID, ok := ctx.Value(dto.RequestIDContextKey).(string)
	if ok {
		return logger.WithField("request_id", requestID)
	}
	return logger
}
