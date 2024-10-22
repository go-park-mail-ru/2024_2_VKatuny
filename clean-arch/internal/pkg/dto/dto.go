package dto

import (
	// "github.com/sirupsen/logrus"
)

type loggerKey int

const LoggerContextKey loggerKey = 1

// Standart json response from backend to frontend
type JsonResponse struct {
	HttpStatus int         `json:"statusCode"`
	Body       interface{} `json:"body"`
	Error      string      `json:"error"`
}
