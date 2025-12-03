package exception

import (
	"net/http"
)

type InternalException struct {
	Messages   []string `json:"messages"`
	StatusText string   `json:"statusText"`
	StatusCode int      `json:"statusCode"`
}

func NewInternalException(messages ...string) *AppException {
	if len(messages) == 0 {
		messages = []string{"erro interno e inesperado"}
	}

	return NewAppException(messages, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
