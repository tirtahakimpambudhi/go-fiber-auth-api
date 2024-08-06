package token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/segmentio/ksuid"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs"
	"time"
)

const MinSecretKeySize = 32

// CustomClaims adds custom payload to JWT claims.
type CustomClaims struct {
	*Payload
	jwt.RegisteredClaims
}

// JWTToken holds configuration for JWT tokens.
type JWTToken struct {
	Name string `env:"TOKEN_NAME,required"`
}

// NewJWTToken initializes a JWTToken with configuration from environment.
func NewJWTToken() (*JWTToken, *SecretKey, error) {
	var (
		jwtToken  JWTToken
		secretKey SecretKey
	)
	err := configs.GetConfig().Load(&jwtToken, &secretKey)
	if err != nil {
		return nil, nil, err
	}
	if len(secretKey.AccessToken) < MinSecretKeySize || len(secretKey.RefreshToken) < MinSecretKeySize || len(secretKey.ForgotPasswordToken) < MinSecretKeySize {
		return nil, nil, NewTokenError(ErrInvalidKey, fmt.Sprintf("min length secret key is %d", MinSecretKeySize))
	}
	return &jwtToken, &secretKey, nil
}

// CreatePayload generates a new payload with specified user ID, email, and expiration.
func (jwtToken JWTToken) CreatePayload(id ksuid.KSUID, email string, duration time.Duration) *Payload {
	return NewTokenPayloadBuilder().WithUserID(id).WithEmail(email).WithExpiration(time.Now().Add(duration)).Build()
}

// CreateToken generates a JWT token with the provided payload.
func (jwtToken JWTToken) CreateToken(secretKey string, payload *Payload) (string, error) {
	// Validate the length of the secret key.
	if len(secretKey) < MinSecretKeySize {
		return "", NewTokenError(ErrInvalidKey, fmt.Sprintf("min length secret key is %d : %d", MinSecretKeySize, len(secretKey)))
	}
	if payload == nil {
		return "", errors.New("token payload cannot be nil")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtToken.Name,
			ExpiresAt: jwt.NewNumericDate(payload.ExpiredAt),
			IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
		},
	})
	return token.SignedString([]byte(secretKey))
}

// VerifyToken parses and verifies a JWT token, returning the payload if valid.
func (jwtToken JWTToken) VerifyToken(secretKey, tokenStr string) (*Payload, error) {
	// Validate the length of the secret key.
	if len(secretKey) < MinSecretKeySize {
		return nil, NewTokenError(ErrInvalidKey, fmt.Sprintf("min length secret key is %d : %d", MinSecretKeySize, len(secretKey)))
	}
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, jwtToken.handleTokenError(err)
	}

	if token.Valid {
		if claims, ok := token.Claims.(*CustomClaims); ok {
			return claims.Payload, nil
		}
		return nil, NewTokenError(ErrFailedParseClaims, "failed to parse claims")
	}

	return nil, NewTokenError(ErrServerError, "invalid token")
}

// handleTokenError provides detailed error messages based on token errors.
func (jwtToken JWTToken) handleTokenError(err error) error {
	switch {
	case errors.Is(err, jwt.ErrTokenMalformed):
		fmt.Println("error malformed token")
		return NewTokenError(ErrTokenMalformed, err.Error())
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		fmt.Println("error signature invalid")
		return NewTokenError(ErrTokenSignatureInvalid, err.Error())
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		fmt.Println("error token expired")
		return NewTokenError(ErrTokenExpired, err.Error())
	default:
		fmt.Println("error internal server")
		return NewTokenError(ErrServerError, err.Error())
	}
}
