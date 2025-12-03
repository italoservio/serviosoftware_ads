package exception

import (
	"net/http"
)

type ResourceExistsException struct {
	Messages   []string `json:"messages"`
	StatusText string   `json:"statusText"`
	StatusCode int      `json:"statusCode"`
}

func NewResourceExistsException(messages ...string) *AppException {
	if len(messages) == 0 {
		messages = []string{"ja existe um recurso com os mesmos atributos"}
	}

	return NewAppException(messages, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
