package exception

import (
	"net/http"
)

type PayloadException struct {
	Messages   []string `json:"messages"`
	StatusText string   `json:"statusText"`
	StatusCode int      `json:"statusCode"`
}

func NewPayloadException(messages ...string) *AppException {
	if len(messages) == 0 {
		messages = []string{"erro inesperado ao processar o contrato da requisicao"}
	}

	return NewAppException(messages, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
