package http

import (
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/bootstrap"
	"github.com/tirtahakimpambudhi/restful_api/internal/repository"
	"github.com/tirtahakimpambudhi/restful_api/internal/usecase"
	"github.com/tirtahakimpambudhi/restful_api/internal/validation"
)

func NewController(app *bootstrap.App) (*UsersController, *AuthController, error) {
	app.Logger.App.Info().Msg("NewController Call Function")
	// Create a new English locale
	english := en.New()

	// Create a new universal translator with English as the default language
	universalTranslate := ut.New(english, english)

	// Get the English translator from the universal translator
	translator, found := universalTranslate.GetTranslator("en")
	if !found {
		app.Logger.App.Error().Msg("the language English not found package")
		return nil, nil, fmt.Errorf("the language English not found package")
	}

	// Create a new UsersRepository implementation
	usersRepository, err := repository.NewUsersRepositoryImpl(app.Gorm, app.Logger.App)
	if err != nil {
		app.Logger.App.Error().Err(err)
		return nil, nil, err
	}

	// Create a new UserCacheRepository instance
	cacheRepository := repository.NewUserCacheRepository(app.Redis.NewClient(), app.Logger.App)

	// Initialize the UsersController with the necessary dependencies
	usersController := NewUsersController(usecase.NewUsersUsecaseBuilder().
		WithHashing(app.Hash).
		WithLogger(app.Logger.App).
		WithUsersRepository(usersRepository).
		WithCacheRepository(cacheRepository).
		WithValidator(validation.NewValidator(validator.New(), translator)).
		WithTimeoutConfig(app.Timeout).
		Build(),
		app.Logger.App)
	// Initialize the AuthController with the necessary dependencies
	authController := NewAuthController(usecase.NewAuthUsecaseBuilder().
		WithHashing(app.Hash).
		WithLogger(app.Logger.App).
		WithUsersRepository(usersRepository).
		WithValidator(validation.NewValidator(validator.New(), translator)).
		WithTimeoutConfig(app.Timeout).
		WithEnforcer(app.CasbinEnforcer).
		WithToken(app.Token).
		WithSecretKey(app.Secret).
		Build(),
		app.Logger.App)
	return usersController, authController, nil
}
