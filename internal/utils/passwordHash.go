package utils

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	cost = 10
)

// HashPassword hashing password with cost og 10
func HashPassword(password string) string {
	bytePassword := []byte(password)
	hashedPassword, _ := bcrypt.GenerateFromPassword(bytePassword, cost)
	return string(hashedPassword[:])
}

// EqualHashedPasswords returns true if first argument was hashed from second
func EqualHashedPasswords(passwordBD string, passwordFront string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordBD), []byte(passwordFront))
	return err == nil
}
