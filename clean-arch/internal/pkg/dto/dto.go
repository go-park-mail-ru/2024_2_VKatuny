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

// use this struct as a field 'Body' in struct JsonResponse
type JsonUserBody struct {
	UserType string `json:"userType"`
	ID       uint64 `json:"id"`
}

