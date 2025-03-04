package api

import "net/http"

type ServiceDeps struct {
}
type Service struct {
}

func NewService(deps *ServiceDeps) *Service {
	return &Service{}
}
func (service *Service) AddCookie(w *http.ResponseWriter, name, value string, maxAge int) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   maxAge,
	}
	http.SetCookie(*w, cookie)
}
