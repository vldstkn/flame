package http_errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

const (
	UserExists            = "user already exists"
	UserDoesNotExist      = "user already exists"
	InvalidNameOrPassword = "invalid email or password"
)

func HandleError(err error) (string, int) {
	st, ok := status.FromError(err)
	if !ok {
		return "", 500
	}
	mes := st.Message()
	var code int
	switch st.Code() {
	case codes.InvalidArgument:
		code = 400
	case codes.Unauthenticated:
		code = 401
	case codes.PermissionDenied:
		code = 403
	case codes.NotFound:
		code = 404
	default:
		code = 500
		mes = http.StatusText(http.StatusInternalServerError)
	}
	return mes, code
}
