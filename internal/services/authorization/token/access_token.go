package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// Генерация токена с полем логина и с определенным временем жизни. Зашифрован ключём.
func GenerateAccessToken(login string, key string, lifeTime time.Duration) (string, error) {
	var claims AuthClaims
	claims.ExpiresAt = time.Now().Add(lifeTime).Unix()
	claims.Login = login

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString([]byte(key))
}

func ParseAccessToken(rawToken string, key string) (string, error) {
	token, parseErr := jwt.ParseWithClaims(rawToken, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(key), nil
	})
	if parseErr != nil {
		return "", parseErr
	}
	claims, ok := token.Claims.(*AuthClaims)
	if !ok {
		return "", errors.New("wrong claims structure")
	} else if claims.ExpiresAt < time.Now().Unix() {
		return "", errors.New("token is expired")
	} else if !token.Valid {
		return "", errors.New("token is not valid")
	}
	return claims.Login, nil
}
