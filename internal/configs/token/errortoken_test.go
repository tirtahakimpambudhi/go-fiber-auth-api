package token_test

import (
	"errors"
	"github.com/stretchr/testify/require"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/token"
	"testing"
)

// Creates a TokenError with a given error and message
func TestCreatesTokenErrorWithGivenErrorAndMessage(t *testing.T) {
	err := token.ErrTokenMalformed
	msg := "This is a test message"
	tokenError := token.NewTokenError(err, msg)

	require.Equal(t, err, tokenError.Err)
	require.True(t, errors.Is(tokenError.TypeError(), err))
	require.Equal(t, msg, tokenError.Message)
}

// Handles nil error input gracefully
func TestHandlesNilErrorInputGracefully(t *testing.T) {
	var err error = nil
	msg := "This is a test message"
	tokenError := token.NewTokenError(err, msg)

	require.Nil(t, tokenError.Err)
	require.Equal(t, msg, tokenError.Message)
}

// Returns formatted error message with both error and message
func TestReturnsFormattedErrorMessage(t *testing.T) {
	err := token.NewTokenError(token.ErrTokenMalformed, "invalid token format")
	expected := "token malformed: invalid token format"
	actual := err.Error()
	require.Equal(t, expected, actual)
}

// Handles empty message string
func TestHandlesEmptyMessageString(t *testing.T) {
	err := token.NewTokenError(token.ErrTokenMalformed, "")
	expected := "token malformed: "
	actual := err.Error()
	require.Equal(t, expected, actual)
}
