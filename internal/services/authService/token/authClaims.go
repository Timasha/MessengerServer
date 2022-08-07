package token

import "github.com/golang-jwt/jwt"

type AuthClaims struct {
	jwt.StandardClaims
	Login string `json:"login"`
}
