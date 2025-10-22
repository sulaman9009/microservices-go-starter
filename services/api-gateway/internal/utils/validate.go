package util

import "github.com/go-playground/validator/v10"

func ValidatePayload(v any) error {
	return validator.New(validator.WithRequiredStructEnabled()).Struct(v)
}
