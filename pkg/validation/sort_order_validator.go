package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func validateSortOrder(fl validator.FieldLevel) bool {
	for _, order := range []string{"asc", "desc"} {
		if strings.EqualFold(fl.Field().String(), order) {
			return true
		}
	}

	return false
}
