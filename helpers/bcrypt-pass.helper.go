package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPass(pass string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
	if err != nil {
		return "", err
	}
	return string(hashedPass), nil
}