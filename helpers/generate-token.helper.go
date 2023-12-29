package helpers

import (
	"os"
	"time"
	"github.com/dgrijalva/jwt-go"
)

var jwtKey []byte

type JWTClaim struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`

	jwt.StandardClaims
}

func init() {
	jwtKey = []byte(os.Getenv("SECRET_KEY"))
	if jwtKey == nil {
		panic("SECRET_KEY not set in .env file")
	}
}

func GenerateToken(id uint, email string) (tokenString string, err error) {
	claims := &JWTClaim{
		ID:    id,
		Email: email,

		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}


