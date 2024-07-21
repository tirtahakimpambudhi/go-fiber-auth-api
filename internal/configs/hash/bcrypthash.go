package hash

import (
	"github.com/tirtahakimpambudhi/restful_api/internal/configs"
	"golang.org/x/crypto/bcrypt"
)

// Bcrypt provides hashing functionality using Bcrypt algorithm
type Bcrypt struct {
	Salt int `env:"HASH_SALT" envDefault:"10"` // Salt for Bcrypt hashing
}

// NewHashBcrypt creates a new Bcrypt instance with configuration from environment
func NewHashBcrypt() (*Bcrypt, error) {
	var hashTech Bcrypt
	// Load configuration and create a new Bcrypt instance
	err := configs.GetConfig().Load(&hashTech)
	if err != nil {
		// Return error if configuration loading fails
		return nil, err
	}
	// Return the Bcrypt instance and error if any
	return &hashTech, err
}

// Create generates a hash for the given password using Bcrypt
func (hash Bcrypt) Create(password string) (string, error) {
	// Generate and return a hash of the password
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), hash.Salt)
	return string(bytes), err
}

// Match checks if the given password matches the stored hash using Bcrypt
func (hash Bcrypt) Match(password, passwordHash string) (bool, error) {
	// Compare the given password with the hash and return if they match
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err == nil, nil
}
