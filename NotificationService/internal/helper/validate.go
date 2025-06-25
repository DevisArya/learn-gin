package helper

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidateStruct(validate *validator.Validate, request interface{}) []string {
	var validationMessages []string

	if err := validate.Struct(request); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range errs {
				validationMessages = append(validationMessages, fmt.Sprintf("%s:%s", e.Field(), e.Tag()))
			}
		}
	}

	return validationMessages
}
