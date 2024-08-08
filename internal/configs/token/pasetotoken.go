package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// PasetoToken holds configuration for PASETO tokens.
type PasetoToken struct {
	paseto *paseto.V2
}

// NewPasetoToken initializes a PasetoToken with configuration from environment.
func NewPasetoToken() (*PasetoToken, *SecretKey, error) {
	var pasetoToken PasetoToken
	secretKey, err := NewSecretKey()
	if err != nil {
		return nil, nil, err
	}
	if len(secretKey.AccessToken) < chacha20poly1305.KeySize || len(secretKey.RefreshToken) < chacha20poly1305.KeySize || len(secretKey.ForgotPasswordToken) < chacha20poly1305.KeySize {
		return nil, nil, NewTokenError(ErrInvalidKey, fmt.Sprintf("min length secret key is %d", MinSecretKeySize))
	}
	pasetoToken.paseto = paseto.NewV2()
	return &pasetoToken, secretKey, nil
}

// CreateToken generates a PASETO token with the provided payload.
func (pasetoToken PasetoToken) CreateToken(secretKey string, payload *Payload) (string, error) {
	if len(secretKey) < chacha20poly1305.KeySize {
		return "", NewTokenError(ErrInvalidKey, fmt.Sprintf("min length secret key is %d : %d", MinSecretKeySize, len(secretKey)))
	}
	if payload == nil {
		return "", errors.New("payload cannot be nil")
	}
	return pasetoToken.paseto.Encrypt([]byte(secretKey), payload, nil)
}

// VerifyToken parses and verifies a PASETO token, returning the payload if valid.
func (pasetoToken PasetoToken) VerifyToken(secretKey string, token string) (*Payload, error) {
	payload := &Payload{}
	if len(secretKey) < chacha20poly1305.KeySize {
		return nil, NewTokenError(ErrInvalidKey, fmt.Sprintf("min length secret key is %d : %d", MinSecretKeySize, len(secretKey)))
	}
	err := pasetoToken.paseto.Decrypt(token, []byte(secretKey), payload, nil)
	if err != nil {
		return nil, NewTokenError(ErrInvalidToken, err.Error())
	}

	if time.Now().After(payload.ExpiredAt) {
		return nil, NewTokenError(ErrExpiredToken, "token is expired")
	}
	return payload, nil
}
