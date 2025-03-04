package mocks

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type MockApiService struct {
	mock.Mock
}

func (mock *MockApiService) AddCookie(w *http.ResponseWriter, name, value string, maxAge int) {
	return
}
