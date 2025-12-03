package exception

import (
	"net/http"
)

type ClientException struct {
	Messages   []string `json:"messages"`
	StatusText string   `json:"statusText"`
	StatusCode int      `json:"statusCode"`
}

func NewClientException(messages ...string) *AppException {
	if len(messages) == 0 {
		messages = []string{"erro inesperado ao processar a requisicao ao cliente"}
	}

	return NewAppException(messages, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
