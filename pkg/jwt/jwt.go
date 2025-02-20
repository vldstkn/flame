package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Data struct {
	Id int64
}

type JWT struct {
	Secret string
}

type Claims struct {
	jwt.RegisteredClaims
	Data
}

func NewJWT(secret string) *JWT {
	return &JWT{
		Secret: secret,
	}
}

func (j *JWT) Create(data Data, expirationTime time.Time) (string, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Data: data,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := t.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}
	return s, nil
}

func (j *JWT) Parse(token string) (bool, *Data) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})

	if err != nil {
		return false, nil
	}

	id := int64(t.Claims.(jwt.MapClaims)["Id"].(float64))
	return t.Valid, &Data{
		Id: id,
	}
}
