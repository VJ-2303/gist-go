package data

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	UserID      int64         `json:"-"`
	TokenString string        `json:"token"`
	Expiry      time.Duration `json:"-"`
}

func GenerateNewToken(token *Token, jwtSecret string) error {
	claims := &jwt.MapClaims{
		"sub": strconv.FormatInt(token.UserID, 10),
		"exp": jwt.NewNumericDate(time.Now().Add(token.Expiry)),
		"iat": jwt.NewNumericDate(time.Now()),
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := tokenClaims.SignedString([]byte(jwtSecret))
	if err != nil {
		return err
	}
	token.TokenString = tokenString

	return nil
}
