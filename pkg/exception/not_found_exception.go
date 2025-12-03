package exception

import (
	"net/http"
)

type NotFoundException struct {
	Messages   []string `json:"messages"`
	StatusText string   `json:"statusText"`
	StatusCode int      `json:"statusCode"`
}

func NewNotFoundException(messages ...string) *AppException {
	if len(messages) == 0 {
		messages = []string{"o recurso solicitado nao foi encontrado"}
	}

	return NewAppException(messages, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
