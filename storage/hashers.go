package storage

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytePassword := []byte(password)
	cost := 10
	hashedPassword, _ := bcrypt.GenerateFromPassword(bytePassword, cost)
	//fmt.Println(string(password[:]), string(hashedPassword[:]))
	return string(hashedPassword[:])
}

func EqualHashedPasswords(passwordBD string, passwordFront string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordBD), []byte(passwordFront))
	if err == nil {
		return true
	} else {
		return false
	}
}
