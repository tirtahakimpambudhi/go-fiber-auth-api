package token

import (
	"errors"
	"github.com/segmentio/ksuid"
	"time"
)

// Define additional error messages specific to PASETO tokens.
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload represents the data contained within a token.
type Payload struct {
	ID        ksuid.KSUID `json:"id"`
	Email     string      `json:"email"`
	IssuedAt  time.Time   `json:"issued_at"`
	ExpiredAt time.Time   `json:"expired_at"`
}

// NewPayload creates a new Payload instance.
func NewPayload(payload *Payload) *Payload {
	return payload
}

// TokenPayloadBuilder helps construct a Payload with a builder pattern.
type TokenPayloadBuilder struct {
	id        ksuid.KSUID
	email     string
	expiredAt time.Time
}

// NewTokenPayloadBuilder creates a new TokenPayloadBuilder instance.
func NewTokenPayloadBuilder() *TokenPayloadBuilder {
	return &TokenPayloadBuilder{}
}

// WithUserID sets the UserID field of the builder.
func (b *TokenPayloadBuilder) WithUserID(id ksuid.KSUID) *TokenPayloadBuilder {
	b.id = id
	return b
}

// WithEmail sets the Email field of the builder.
func (b *TokenPayloadBuilder) WithEmail(email string) *TokenPayloadBuilder {
	b.email = email
	return b
}

// WithExpiration sets the Expiration field of the builder.
func (b *TokenPayloadBuilder) WithExpiration(expired time.Time) *TokenPayloadBuilder {
	b.expiredAt = expired
	return b
}

// Build creates a Payload from the builder.
func (b *TokenPayloadBuilder) Build() *Payload {
	return &Payload{
		ID:        b.id,
		Email:     b.email,
		IssuedAt:  time.Now(),
		ExpiredAt: b.expiredAt,
	}
}
