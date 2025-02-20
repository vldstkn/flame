package interfaces

import "net/http"

type ApiService interface {
	AddCookie(w *http.ResponseWriter, name, value string, maxAge int)
}
