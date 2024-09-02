package handler

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	logrus.Error(message)

	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}
