package hash

import (
	"github.com/alexedwards/argon2id"
)

// Argon2 provides hashing functionality using Argon2 algorithm
type Argon2 struct {
	Params *argon2id.Params // Parameters for Argon2 hashing
}

// NewHashArgon2 creates a new Argon2 instance with default parameters
func NewHashArgon2() (*Argon2, error) {
	// Return a new Argon2 instance with default parameters
	return &Argon2{
		Params: argon2id.DefaultParams,
	}, nil
}

// Create generates a hash for the given password using Argon2
func (hash Argon2) Create(password string) (passwordHash string, err error) {
	// Create and return a hash of the password
	return argon2id.CreateHash(password, hash.Params)
}

// Match checks if the given password matches the stored hash using Argon2
func (hash Argon2) Match(password string, passwordHash string) (match bool, err error) {
	// Compare the given password with the hash and return if they match
	return argon2id.ComparePasswordAndHash(password, passwordHash)
}
