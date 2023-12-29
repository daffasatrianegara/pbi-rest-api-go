package helpers

import (
	"errors"
	"time"
	"github.com/dgrijalva/jwt-go"
)

func ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("token invalid")
	}

	if token.Claims.(*JWTClaim).ExpiresAt < time.Now().Unix() {
		return errors.New("token expired")
	}

	return nil
}