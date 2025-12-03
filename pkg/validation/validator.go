package validation

import "github.com/go-playground/validator/v10"

func NewValidator() (*validator.Validate, error) {
	validate := validator.New()

	if err := validate.RegisterValidation("oneofsortorder", validateSortOrder); err != nil {
		return nil, err
	}

	return validate, nil

}
