package transport

import (
	"ride-sharing/services/api-gateway/internal/problems"

	"github.com/go-playground/validator/v10"
)

type customValidator struct {
	validator *validator.Validate
}

func NewCustomValidator() *customValidator {
	return &customValidator{
		validator: validator.New(),
	}
}

func (cv *customValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return problems.NewBadRequest(err.Error(), "")
	}
	return nil
}
