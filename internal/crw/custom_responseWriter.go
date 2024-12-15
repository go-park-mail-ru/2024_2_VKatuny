package crw

import (
	"net/http"
)

type ResponseWriterStatus struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriterStatus(w http.ResponseWriter) *ResponseWriterStatus {
	ws := new(ResponseWriterStatus)
	ws.ResponseWriter = w
	return ws
}

func (ws *ResponseWriterStatus) WriteHeader(statusCode int) {
	ws.statusCode = statusCode
	ws.ResponseWriter.WriteHeader(statusCode)
}

func (ws *ResponseWriterStatus) Status() int {
	return ws.statusCode
}
