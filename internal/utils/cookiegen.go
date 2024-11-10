package utils

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

const TokenLength = 32

// GenerateSessionTokenWithLength generate random string with given length for session id
func GenerateSessionToken(n int, userType string) string {
	b := make([]rune, n)

	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	if userType == dto.UserTypeApplicant {
		return "1" + string(b)
	}
	return "2" + string(b)
}

func CheckToken(token string) (string, error) {
	if token[0] == '1' {
		return dto.UserTypeApplicant, nil
	} else if token[0] == '2' {
		return dto.UserTypeEmployer, nil
	} else {
		return "", fmt.Errorf("not our cookie")
	}
}

func MakeAuthCookie(sessionID, backendAddress string) *http.Cookie {
	return &http.Cookie{
		Name:     "session_id1", // why id1?
		Value:    sessionID,
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
		//Secure:   true, // Enable when using HTTPS
		SameSite: http.SameSiteStrictMode,
		Domain:   backendAddress,
		Path:     "/api/v1",
	}
}
