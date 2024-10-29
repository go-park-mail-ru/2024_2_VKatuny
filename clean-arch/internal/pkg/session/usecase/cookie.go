package usecase

import (
	"math/rand"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

const tokenLength = 32

func GenerateSessionToken() string {
	return GenerateSessionTokenWithLength(tokenLength)
}

// RandStringRunes generate random string with length of n for session id
func GenerateSessionTokenWithLength(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
