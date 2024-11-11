package utils

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	"github.com/sirupsen/logrus"
)

func SetRequestIDInLoggerFromRequest(r *http.Request, logger *logrus.Entry) *logrus.Entry {
	requestID, ok := r.Context().Value(dto.RequestIDContextKey).(string)
	if ok {
		return logger.WithField("request_id", requestID)
	} 
	logger.Errorf("unable to get requestID from context")
	return logger
}
