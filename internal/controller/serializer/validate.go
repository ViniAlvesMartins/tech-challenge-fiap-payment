package serializer

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type ValidateError struct {
	Errors []Fields
}

type Fields struct {
	Field   string
	Message string
}

func Validate(dto interface{}) ValidateError {
	var validateError ValidateError
	var errList []Fields

	validate = validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(dto); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errList = append(errList, Fields{
				Field:   e.Field(),
				Message: e.Tag(),
			})
		}
	}

	validateError.Errors = errList

	return validateError
}
