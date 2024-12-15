package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

const TokenExpirationTime = 30 * time.Minute

var (
	ErrShortCipher      = fmt.Errorf("short cipher")
	ErrDecryptionFailed = fmt.Errorf("decryption failed")
	ErrTokenExpired     = fmt.Errorf("csrf token expired")
	ErrInvalidKeyLength = fmt.Errorf("aes key length must be 16, 24 or 32 bytes")
)

type CryptToken struct {
	secret []byte
}

type CSRFTokenData struct {
	SessionID string
	UserID    uint64
	UserType  string
	TTL       int64
}

func NewCryptToken(secret string) (*CryptToken, error) {
	key := []byte(secret)
	keyLength := len(key)
	if keyLength != 16 && keyLength != 24 && keyLength != 32 {
		return nil, ErrInvalidKeyLength
	}
	return &CryptToken{secret: key}, nil
}

func (c *CryptToken) Create(userID uint64, UserType, SessionID string) (string, error) {
	block, err := aes.NewCipher(c.secret)
	if err != nil {
		return "", err
	}

	AESGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, AESGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	Token := &CSRFTokenData{
		SessionID: SessionID,
		UserID:    userID,
		UserType:  UserType,
		TTL:       time.Now().Add(TokenExpirationTime).Unix(),
	}
	TokenJSON, err := json.Marshal(Token)
	if err != nil {
		return "", err
	}
	cipherText := AESGCM.Seal(nil, nonce, TokenJSON, nil)

	result := make([]byte, 0)
	result = append(result, nonce...)
	result = append(result, cipherText...)

	return base64.StdEncoding.EncodeToString(result), nil
}

func (c *CryptToken) Check(userID uint64, UserType, SessionID, inputToken string) (bool, error) {
	token, err := base64.StdEncoding.DecodeString(inputToken)
	if err != nil {
		return false, err
	}

	block, err := aes.NewCipher(c.secret)
	if err != nil {
		return false, err
	}
	AESGCM, err := cipher.NewGCM(block)
	if err != nil {
		return false, err
	}

	nonceSize := AESGCM.NonceSize()
	if len(token) < nonceSize {
		return false, ErrShortCipher
	}

	nonce, cipherText := token[:nonceSize], token[nonceSize:]
	decodedText, err := AESGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return false, ErrDecryptionFailed
	}

	tokenData := new(CSRFTokenData)
	err = json.Unmarshal(decodedText, tokenData)
	if err != nil {
		return false, err
	}

	if tokenData.TTL < time.Now().Unix() {
		return false, ErrTokenExpired
	}

	return tokenData.UserID == userID &&
		tokenData.SessionID == SessionID &&
		tokenData.UserType == UserType, nil
}
