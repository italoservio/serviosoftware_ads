package exception

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type AppException struct {
	Messages   []string `json:"messages"`
	StatusText string   `json:"statusText"`
	StatusCode int      `json:"statusCode"`
}

func NewAppException(messages []string, statusText string, statusCode int) *AppException {
	return &AppException{
		Messages:   messages,
		StatusText: statusText,
		StatusCode: statusCode,
	}
}

func ToAppException(err error) *AppException {
	var appException *AppException

	if errors.As(err, &appException) {
		return appException
	}

	return NewAppException(
		[]string{err.Error()},
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}

func (e *AppException) Error() string {
	return strings.Join(e.Messages, "; ")
}

func (e *AppException) WriteJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.StatusCode)

	json.NewEncoder(w).Encode(e)
}
