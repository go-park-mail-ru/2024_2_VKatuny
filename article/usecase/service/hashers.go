// Package service is business logic layer
package service

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashing password with cost og 10
func HashPassword(password string) string {
	bytePassword := []byte(password)
	cost := 10
	hashedPassword, _ := bcrypt.GenerateFromPassword(bytePassword, cost)
	return string(hashedPassword[:])
}

// EqualHashedPasswords returns true if first argument was hashed from second
func EqualHashedPasswords(passwordBD string, passwordFront string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordBD), []byte(passwordFront))
	return err == nil
}
