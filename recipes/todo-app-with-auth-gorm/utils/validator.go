package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"go.khulnasoft.com/velocity"
)

var validate = validator.New()

// Validate validates the input struct
func Validate(payload interface{}) *velocity.Error {
	err := validate.Struct(payload)
	if err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(
				errors,
				fmt.Sprintf("`%v` with value `%v` doesn't satisfy the `%v` constraint", err.Field(), err.Value(), err.Tag()),
			)
		}

		return &velocity.Error{
			Code:    velocity.StatusBadRequest,
			Message: strings.Join(errors, ","),
		}
	}

	return nil
}

// CUSTOM VALIDATION RULES =============================================

// Password validation rule: required,min=6,max=100
var _ = validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
	l := len(fl.Field().String())

	return l >= 6 && l < 100
})
