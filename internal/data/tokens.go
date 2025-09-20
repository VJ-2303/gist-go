package data

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	UserID      int64         `json:"-"`
	TokenString string        `json:"token"`
	Expiry      time.Duration `json:"-"`
}

func GenerateNewToken(token *Token, jwtSecret string) error {
	claims := &jwt.RegisteredClaims{
		Subject:   string(rune(token.UserID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(token.Expiry)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := tokenClaims.SignedString([]byte(jwtSecret))
	if err != nil {
		return err
	}
	token.TokenString = tokenString

	return nil
}
