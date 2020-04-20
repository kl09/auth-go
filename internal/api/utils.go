package api

import (
	"golang.org/x/crypto/bcrypt"
)

func hashAndSalt(pwd string) (hash string, err error) {
	h, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(h), err
}

func comparePasswords(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plainPwd))

	return err == nil
}
