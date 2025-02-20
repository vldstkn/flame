package req

import "github.com/go-playground/validator/v10"

func IsValid[T any](payload T) error {
	err := validator.New().Struct(payload)
	return err
}
