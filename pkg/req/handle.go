package req

import (
	"net/http"
)

func HandleBody[T any](r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)
	if err != nil {
		return nil, err
	}
	err = IsValid(body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}
