package usecase

import (
	"github.com/casbin/casbin/v2"
	"github.com/phuslu/log"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/hash"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/timeout"
	tokenconfig "github.com/tirtahakimpambudhi/restful_api/internal/configs/token"
	"github.com/tirtahakimpambudhi/restful_api/internal/entity"
	"github.com/tirtahakimpambudhi/restful_api/internal/repository"
	"github.com/tirtahakimpambudhi/restful_api/internal/validation"
)

// AuthUsecaseBuilder is the builder for AuthUsecase.
type AuthUsecaseBuilder struct {
	usersRepository repository.UsersRepository
	timeoutConfig   *timeout.Config
	validator       *validation.Validator
	hashing         *hash.Argon2
	token           *tokenconfig.JWTToken
	secretKey       *tokenconfig.SecretKey
	logger          *log.Logger
	enforcer        *casbin.Enforcer
}

// NewAuthUsecaseBuilder creates a new instance of AuthUsecaseBuilder.
func NewAuthUsecaseBuilder() *AuthUsecaseBuilder {
	return &AuthUsecaseBuilder{}
}

// WithUsersRepository sets the UsersRepository.
func (b *AuthUsecaseBuilder) WithUsersRepository(repo repository.UsersRepository) *AuthUsecaseBuilder {
	b.usersRepository = repo
	return b
}

// WithTimeoutConfig sets the TimeoutConfig.
func (b *AuthUsecaseBuilder) WithTimeoutConfig(timeout *timeout.Config) *AuthUsecaseBuilder {
	b.timeoutConfig = timeout
	return b
}

// WithValidator sets the Validator.
func (b *AuthUsecaseBuilder) WithValidator(validator *validation.Validator) *AuthUsecaseBuilder {
	b.validator = validator
	return b
}

// WithHashing sets the Argon2 hashing utility.
func (b *AuthUsecaseBuilder) WithHashing(hashing *hash.Argon2) *AuthUsecaseBuilder {
	b.hashing = hashing
	return b
}

// WithToken sets the PasetoToken utility.
func (b *AuthUsecaseBuilder) WithToken(token *tokenconfig.JWTToken) *AuthUsecaseBuilder {
	b.token = token
	return b
}

// WithSecretKey sets the SecretKey.
func (b *AuthUsecaseBuilder) WithSecretKey(secretKey *tokenconfig.SecretKey) *AuthUsecaseBuilder {
	b.secretKey = secretKey
	return b
}

// WithLogger sets the Logger.
func (b *AuthUsecaseBuilder) WithLogger(logger *log.Logger) *AuthUsecaseBuilder {
	b.logger = logger
	return b
}

// WithEnforcer sets the enforcer
func (b *AuthUsecaseBuilder) WithEnforcer(enforcer *casbin.Enforcer) *AuthUsecaseBuilder {
	b.enforcer = enforcer
	return b
}

// Build creates the AuthUsecase instance.
func (b *AuthUsecaseBuilder) Build() *AuthUsecase {
	return &AuthUsecase{
		usersRepository: b.usersRepository,
		timeoutConfig:   b.timeoutConfig,
		validator:       b.validator,
		hashing:         b.hashing,
		token:           b.token,
		secretKey:       b.secretKey,
		logger:          b.logger,
	}
}

// UsersUsecaseBuilder is the builder for UsersUsecase.
type UsersUsecaseBuilder struct {
	usersRepository repository.UsersRepository
	cacheRepository repository.CacheRepository[*entity.Users]
	timeoutConfig   *timeout.Config
	validator       *validation.Validator
	hashing         *hash.Argon2
	logger          *log.Logger
}

// NewUsersUsecaseBuilder creates a new instance of UsersUsecaseBuilder.
func NewUsersUsecaseBuilder() *UsersUsecaseBuilder {
	return &UsersUsecaseBuilder{}
}

// WithUsersRepository sets the UsersRepository.
func (b *UsersUsecaseBuilder) WithUsersRepository(repo repository.UsersRepository) *UsersUsecaseBuilder {
	b.usersRepository = repo
	return b
}

// WithCacheRepository sets the CacheRepository.
func (b *UsersUsecaseBuilder) WithCacheRepository(cacheRepo repository.CacheRepository[*entity.Users]) *UsersUsecaseBuilder {
	b.cacheRepository = cacheRepo
	return b
}

// WithTimeoutConfig sets the TimeoutConfig.
func (b *UsersUsecaseBuilder) WithTimeoutConfig(timeout *timeout.Config) *UsersUsecaseBuilder {
	b.timeoutConfig = timeout
	return b
}

// WithValidator sets the Validator.
func (b *UsersUsecaseBuilder) WithValidator(validator *validation.Validator) *UsersUsecaseBuilder {
	b.validator = validator
	return b
}

// WithHashing sets the Argon2 hashing utility.
func (b *UsersUsecaseBuilder) WithHashing(hashing *hash.Argon2) *UsersUsecaseBuilder {
	b.hashing = hashing
	return b
}

// WithLogger sets the Logger.
func (b *UsersUsecaseBuilder) WithLogger(logger *log.Logger) *UsersUsecaseBuilder {
	b.logger = logger
	return b
}

// Build creates the UsersUsecase instance.
func (b *UsersUsecaseBuilder) Build() *UsersUsecase {
	return &UsersUsecase{
		usersRepository: b.usersRepository,
		cacheRepository: b.cacheRepository,
		timeoutConfig:   b.timeoutConfig,
		validator:       b.validator,
		hashing:         b.hashing,
		logger:          b.logger,
	}
}
