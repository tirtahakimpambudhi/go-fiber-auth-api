package token_test

import (
	"errors"
	"fmt"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
	token "github.com/tirtahakimpambudhi/restful_api/internal/configs/token"
	"os"
	"testing"
	"time"
)

func NewTestJWTToken(t *testing.T) (*token.JWTToken, func()) {
	//	set and unset environment variable
	os.Setenv("TOKEN_NAME", "testing_jwt_token")
	os.Setenv("SECRET_KEY_ACCESS_TOKEN", "a_very_secret_key_that_is_32_byt")
	os.Setenv("SECRET_KEY_REFRESH_TOKEN", "a_very_secret_key_that_is_32_byt")
	os.Setenv("SECRET_KEY_FP_TOKEN", "a_very_secret_key_that_is_32_byt")
	unset := func() {
		os.Unsetenv("SECRET_KEY_ACCESS_TOKEN")
		os.Unsetenv("SECRET_KEY_REFRESH_TOKEN")
		os.Unsetenv("SECRET_KEY_FP_TOKEN")
	}
	jwtToken, _, err := token.NewJWTToken()
	require.NoError(t, err)
	return jwtToken, unset
}

// Successfully loads configuration into JWTToken struct
func TestNewJWTTokenSuccess(t *testing.T) {
	testCases := []struct {
		name        string
		tokenName   string
		tokenSecret string
		isErr       bool
	}{
		{name: "Successfully Initialize JWT Token", tokenName: "testing_jwt_token", tokenSecret: "a_very_secret_key_that_is_long_enough", isErr: false},
		{name: "Failure Initialize JWT Token Because SecretKey to short", tokenName: "", tokenSecret: "", isErr: true},
	}
	for index, testCase := range testCases {
		nameCase := fmt.Sprintf("Case #%d: %s", index, testCase.name)
		t.Run(nameCase, func(t *testing.T) {
			//	set and unset environment variable
			os.Setenv("TOKEN_NAME", testCase.tokenName)
			os.Setenv("SECRET_KEY_ACCESS_TOKEN", testCase.tokenSecret)
			os.Setenv("SECRET_KEY_REFRESH_TOKEN", testCase.tokenSecret)
			os.Setenv("SECRET_KEY_FP_TOKEN", testCase.tokenSecret)
			defer func() {
				os.Unsetenv("SECRET_KEY_ACCESS_TOKEN")
				os.Unsetenv("SECRET_KEY_REFRESH_TOKEN")
				os.Unsetenv("SECRET_KEY_FP_TOKEN")
			}()

			// initialize jwt token
			jwtToken, secretKey, err := token.NewJWTToken()

			// Assert
			if testCase.isErr {
				require.Error(t, err)
				return
			}
			require.Equal(t, testCase.tokenName, jwtToken.Name)
			require.Equal(t, testCase.tokenSecret, secretKey.AccessToken)
			require.Equal(t, testCase.tokenSecret, secretKey.RefreshToken)
			require.Equal(t, testCase.tokenSecret, secretKey.ForgotPasswordToken)
		})
	}
}

// CreatePayload generates a Payload with valid UUID, email, role, and duration
func TestCreatePayloadWithValidData(t *testing.T) {
	//	set and unset environment variable
	os.Setenv("TOKEN_NAME", "testing_jwt_token")
	os.Setenv("SECRET_KEY_ACCESS_TOKEN", "a_very_secret_key_that_is_32_byt")
	os.Setenv("SECRET_KEY_REFRESH_TOKEN", "a_very_secret_key_that_is_32_byt")
	os.Setenv("SECRET_KEY_FP_TOKEN", "a_very_secret_key_that_is_32_byt")
	defer func() {
		os.Unsetenv("SECRET_KEY_ACCESS_TOKEN")
		os.Unsetenv("SECRET_KEY_REFRESH_TOKEN")
		os.Unsetenv("SECRET_KEY_FP_TOKEN")
	}()
	jwtToken, _, err := token.NewJWTToken()

	id := ksuid.New()
	email := "test@example.com"
	duration := time.Hour

	require.NoError(t, err)
	payload := jwtToken.CreatePayload(id, email, duration)
	require.NotNil(t, payload)
	require.Equal(t, id, payload.ID)
	require.Equal(t, email, payload.Email)
	require.True(t, payload.ExpiredAt.After(time.Now()))
}

// Handles the case where the payload is nil
func TestCreateTokenWithNilPayload(t *testing.T) {
	jwtToken, unset := NewTestJWTToken(t)
	defer unset()

	tokenString, err := jwtToken.CreateToken("", nil)
	require.Error(t, err)
	require.Equal(t, "", tokenString)
}

// VerifyToken returns payload for valid token
func TestVerifyTokenReturnsPayloadForValidToken(t *testing.T) {
	os.Setenv("TOKEN_NAME", "testing_jwt_token")
	os.Setenv("SECRET_KEY_ACCESS_TOKEN", "a_very_secret_key_that_is_32_byt")
	os.Setenv("SECRET_KEY_REFRESH_TOKEN", "a_very_secret_key_that_is_32_byt")
	os.Setenv("SECRET_KEY_FP_TOKEN", "a_very_secret_key_that_is_32_byt")
	defer func() {
		os.Unsetenv("SECRET_KEY_ACCESS_TOKEN")
		os.Unsetenv("SECRET_KEY_REFRESH_TOKEN")
		os.Unsetenv("SECRET_KEY_FP_TOKEN")
	}()
	jwtToken, secretKey, err := token.NewJWTToken()
	require.NoError(t, err)

	payload := &token.Payload{
		ID:        ksuid.New(),
		Email:     "test@example.com",
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(time.Hour),
	}

	tokenStr, err := jwtToken.CreateToken(secretKey.AccessToken, payload)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	returnedPayload, err := jwtToken.VerifyToken(secretKey.AccessToken, tokenStr)
	require.NoError(t, err)

	require.Equal(t, payload.ID, returnedPayload.ID)
}

// VerifyToken handles malformed token error
func TestVerifyTokenHandlesMalformedTokenError(t *testing.T) {
	os.Setenv("TOKEN_NAME", "testing_jwt_token")
	os.Setenv("SECRET_KEY_ACCESS_TOKEN", "a_very_secret_key_that_is_32_byt")
	os.Setenv("SECRET_KEY_REFRESH_TOKEN", "a_very_secret_key_that_is_32_byt")
	os.Setenv("SECRET_KEY_FP_TOKEN", "a_very_secret_key_that_is_32_byt")
	defer func() {
		os.Unsetenv("SECRET_KEY_ACCESS_TOKEN")
		os.Unsetenv("SECRET_KEY_REFRESH_TOKEN")
		os.Unsetenv("SECRET_KEY_FP_TOKEN")
	}()
	jwtToken, secretKey, err := token.NewJWTToken()
	require.NoError(t, err)

	malformedToken := "this.is.a.malformed.token"

	_, err = jwtToken.VerifyToken(secretKey.AccessToken, malformedToken)
	if err == nil {
		t.Fatal("Expected an error for malformed token, got nil")
	}

	var tokenErr *token.TokenError
	if !errors.As(err, &tokenErr) || !errors.Is(tokenErr.Err, token.ErrTokenMalformed) {
		t.Errorf("Expected ErrTokenMalformed, got %v", err)
	}
}

func TestVerifyTokenHandlesExpiredToken(t *testing.T) {
	os.Setenv("TOKEN_NAME", "testing_jwt_token")
	os.Setenv("SECRET_KEY_ACCESS_TOKEN", "a_very_secret_key_that_is_32_byt")
	os.Setenv("SECRET_KEY_REFRESH_TOKEN", "a_very_secret_key_that_is_32_byt")
	os.Setenv("SECRET_KEY_FP_TOKEN", "a_very_secret_key_that_is_32_byt")
	defer func() {
		os.Unsetenv("SECRET_KEY_ACCESS_TOKEN")
		os.Unsetenv("SECRET_KEY_REFRESH_TOKEN")
		os.Unsetenv("SECRET_KEY_FP_TOKEN")
	}()
	jwtToken, secretKey, err := token.NewJWTToken()
	require.NoError(t, err)

	tokenStr, err := jwtToken.CreateToken(secretKey.AccessToken, new(token.Payload))
	require.NoError(t, err)

	_, err = jwtToken.VerifyToken(secretKey.AccessToken, tokenStr)
	var tokenErr *token.TokenError
	if !errors.As(err, &tokenErr) || !errors.Is(tokenErr.Err, token.ErrTokenExpired) {
		t.Errorf("Expected ErrTokenMalformed, got %v", err)
	}
}
