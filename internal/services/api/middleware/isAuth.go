package middleware

import (
	"context"
	"flame/internal/services/api/dto"
	"flame/pkg/jwt"
	"flame/pkg/res"
	"net/http"
	"strings"
)

type AuthData struct {
	Id int64
}

func writeUnauthed(w http.ResponseWriter) {
	res.Json(w, dto.ErrorRes{
		Error: http.StatusText(http.StatusUnauthorized),
	}, http.StatusUnauthorized)
}

func IsAuthed(secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authedHeader := r.Header.Get("Authorization")
			if authedHeader == "" || !strings.HasPrefix(authedHeader, "Bearer ") {
				writeUnauthed(w)
				return
			}
			token := strings.TrimPrefix(authedHeader, "Bearer ")
			isValid, data := jwt.NewJWT(secret).Parse(token)

			if !isValid {
				writeUnauthed(w)
				return
			}
			ctx := context.WithValue(r.Context(), "authData", AuthData{
				Id: data.Id,
			})
			req := r.WithContext(ctx)
			next.ServeHTTP(w, req)
		})
	}
}
