package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(key []byte, login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		Subject:   login,
	})
	return token.SignedString(key)
}
func ParseToken(tokenString, key string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected singing method: %v", t.Method.Alg())
		}
		return []byte(key), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims["sub"].(string), nil
	} else {
		return "", fmt.Errorf("Error cast to jwt.MapClaims")
	}
}
