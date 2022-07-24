package token

import (
	"time"

	"github.com/golang-jwt/jwt"
)

// Генерация токена с полем логина и с определенным временем жизни. Зашифрован ключём.
func GenerateAccessJWT(login string, key string, lifeTime time.Duration) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(lifeTime).Unix()
	claims["login"] = login

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString([]byte(key))
}
