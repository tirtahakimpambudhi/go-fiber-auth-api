package token_test

import (
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
	token "github.com/tirtahakimpambudhi/restful_api/internal/configs/token"
	"log"
	"os"
	"testing"
	"time"
)

func NewTestPaseto() (*token.PasetoToken, *token.SecretKey, func()) {
	os.Setenv("SECRET_KEY_ACCESS_TOKEN", "a_very_secret_key_that_is_32_byt")
	os.Setenv("SECRET_KEY_REFRESH_TOKEN", "a_very_secret_key_that_is_32_byt")
	os.Setenv("SECRET_KEY_FP_TOKEN", "a_very_secret_key_that_is_32_byt")
	unsetFunc := func() {
		os.Unsetenv("SECRET_KEY_ACCESS_TOKEN")
		os.Unsetenv("SECRET_KEY_REFRESH_TOKEN")
		os.Unsetenv("SECRET_KEY_FP_TOKEN")
	}

	// Execute
	pasetoToken, secretKey, err := token.NewPasetoToken()
	if err != nil {
		log.Printf("Error creating token : %s \n", err.Error())
		return nil, nil, unsetFunc
	}
	return pasetoToken, secretKey, unsetFunc
}

// Successfully loads configuration into PasetoToken struct
func TestNewPasetoToken_Success(t *testing.T) {
	// Setup
	pasetoToken, secretKey, unsetFunc := NewTestPaseto()
	defer unsetFunc()
	// Verify
	require.NotNil(t, pasetoToken)
	require.Equal(t, "a_very_secret_key_that_is_32_byt", secretKey.AccessToken)
}

// Failure loads configuration into PasetoToken struct
func TestNewPasetoToken_Failure(t *testing.T) {
	// Setup
	os.Setenv("TOKEN_SECRET", "a_very_secret_key_that_is_32_byte")
	defer os.Unsetenv("TOKEN_SECRET")

	// Execute
	token, _, err := token.NewPasetoToken()

	// Verify
	require.Error(t, err)
	require.Nil(t, token)
}

// Successfully creates a token with valid payload and symmetric key
func TestCreateTokenWithValidPayload(t *testing.T) {
	// Setup
	pasetoToken, secretKey, unsetFunc := NewTestPaseto()
	defer unsetFunc()

	payload := token.NewTokenPayloadBuilder().WithUserID(ksuid.New()).WithEmail("test@gmail.com").WithExpiration(time.Now().Add(24 * time.Hour)).Build()
	// Act
	tokenString, err := pasetoToken.CreateToken(secretKey.AccessToken, payload)

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)
}

// Handles empty payload gracefully
func TestCreateTokenWithEmptyPayload(t *testing.T) {
	// Initialize PasetoToken with a valid symmetric key
	// Setup
	pasetoToken, secretKey, unsetFunc := NewTestPaseto()
	defer unsetFunc()

	// Create an empty payload
	var payload *token.Payload

	// Call CreateToken
	tokenString, err := pasetoToken.CreateToken(secretKey.AccessToken, payload)

	require.Error(t, err)
	require.Empty(t, tokenString)
}

// Successfully decrypts a valid token
func TestVerifyToken_Success(t *testing.T) {
	// Setup
	pasetoToken, secretKey, unsetFunc := NewTestPaseto()
	defer unsetFunc()

	payload := token.NewTokenPayloadBuilder().WithUserID(ksuid.New()).WithEmail("test@gmail.com").WithExpiration(time.Now().Add(24 * time.Hour)).Build()
	tokenString, err := pasetoToken.CreateToken(secretKey.AccessToken, payload)

	// Act
	result, errVerify := pasetoToken.VerifyToken(secretKey.AccessToken, tokenString)

	// Assert
	require.NoError(t, err)
	require.NoError(t, errVerify)
	require.NotNil(t, result)
}

func TestVerifyToken_Failure(t *testing.T) {
	// Setup
	pasetoToken, secretKey, unsetFunc := NewTestPaseto()
	defer unsetFunc()

	tokenString := ""
	// Act
	result, err := pasetoToken.VerifyToken(secretKey.AccessToken, tokenString)

	// Assert
	require.Error(t, err)
	require.Nil(t, result)
}
