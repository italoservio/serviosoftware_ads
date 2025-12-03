package exception

import (
	"net/http"
)

type UnauthorizedException struct {
	Messages   []string `json:"messages"`
	StatusText string   `json:"statusText"`
	StatusCode int      `json:"statusCode"`
}

func NewUnauthorizedException(messages ...string) *AppException {
	if len(messages) == 0 {
		messages = []string{"acesso ao recurso nao autorizado"}
	}

	return NewAppException(messages, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}
