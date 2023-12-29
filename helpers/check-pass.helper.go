package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

func CheckPass(pass, hash string)(bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		return false, err
	}
	return true, nil
}