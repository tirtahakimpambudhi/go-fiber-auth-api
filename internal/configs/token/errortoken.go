package token

import (
	"errors"
	"fmt"
)

// Define common error messages for token operations.
var (
	ErrTokenMalformed        = errors.New("token malformed")
	ErrTokenSignatureInvalid = errors.New("token signature invalid")
	ErrTokenExpired          = errors.New("token expired or not valid yet")
	ErrFailedParseClaims     = errors.New("failed to parse claims")
	ErrFailedParseID         = errors.New("failed to parse id")
	ErrServerError           = errors.New("server error")
	ErrInvalidKey            = errors.New("invalid secret key")
)

// TokenError represents an error with additional context.
type TokenError struct {
	Err     error
	Message string
}

// Error returns the error message in a formatted string.
func (e *TokenError) Error() string {
	return fmt.Sprintf("%v: %s", e.Err, e.Message)
}

// TypeError returns the type error in a formatted string.
func (e *TokenError) TypeError() error {
	return e.Err
}

// NewTokenError creates a new TokenError with the provided error and message.
func NewTokenError(err error, msg string) *TokenError {
	return &TokenError{
		Err:     err,
		Message: msg,
		}
	}