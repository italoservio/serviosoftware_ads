package exception

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidatorException struct {
	Messages   []string `json:"messages"`
	StatusText string   `json:"statusText"`
	StatusCode int      `json:"statusCode"`
}

func NewValidatorException(err error) *AppException {
	var validationErrors validator.ValidationErrors

	if !errors.As(err, &validationErrors) {
		return NewAppException(
			[]string{err.Error()},
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
		)
	}

	var messages = []string{"um ou mais campos falharam na validacao"}

	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		tag := err.Tag()

		if len(field) > 0 {
			field = strings.ToLower(string(field[0])) + field[1:]
		}

		messages = append(messages, fmt.Sprintf("o campo '%s' falhou na validacao '%s'", field, tag))
	}

	return NewAppException(messages, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
