package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/mapper"
	"math"
	"net/http"

	"github.com/phuslu/log"
	"github.com/segmentio/ksuid"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/hash"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/timeout"
	"github.com/tirtahakimpambudhi/restful_api/internal/entity"
	errorshandler "github.com/tirtahakimpambudhi/restful_api/internal/errors"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/request"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/response"
	"github.com/tirtahakimpambudhi/restful_api/internal/repository"
	"github.com/tirtahakimpambudhi/restful_api/internal/validation"
)

// UsersUsecase represents the use case layer that handles business logic
// related to users, including repository interactions, caching, and validation.
type UsersUsecase struct {
	usersRepository repository.UsersRepository
	cacheRepository repository.CacheRepository[*entity.Users]
	timeoutConfig   *timeout.Config
	validator       *validation.Validator
	hashing         *hash.Argon2
	logger          *log.Logger
}

// List retrieves a list of users based on the provided request parameters.
func (usersUsecase UsersUsecase) List(ctx context.Context, request *request.Page) (*response.LinksAble, *response.StandardErrors) {
	usersUsecase.logger.Info().Msg("List method called")

	// Validate the incoming request parameters.
	if errValidate := usersUsecase.validator.Validate(request); errValidate != nil {
		usersUsecase.logger.Error().Msgf("Validation error: %v", errValidate)
		return nil, &response.StandardErrors{Errors: errValidate}
	}
	usersUsecase.logger.Info().Msg("Request validated successfully")

	// Set a timeout context for cache operations.
	ctxTimeout, cancel := usersUsecase.timeoutConfig.CreateCacheTimeout(ctx)
	defer cancel() // Ensure context cancellation after function completes.

	// Generate a cache key based on request parameters.
	key := fmt.Sprintf("users:all:size[%d]", request.Size)
	if request.Before != "" {
		key += fmt.Sprintf(":before[%s]", request.Before)
	}
	if request.After != "" {
		key += fmt.Sprintf(":after[%s]", request.After)
	}
	usersUsecase.logger.Info().Msgf("Cache key generated: %s", key)

	// Attempt to retrieve the users list from the cache.
	users, errCache := usersUsecase.cacheRepository.GetFromCache(ctxTimeout, key)
	usersUsecase.logger.Info().Msgf("Errors: %v", errCache)
	if errCache != nil {
		usersUsecase.logger.Error().Msgf("Failed to fetch users from cache: %v", errCache)
		return nil, usersUsecase.handleErrFromRepository(errCache, "Failed to fetch users from cache: ")
	}
	if users == nil {
		usersUsecase.logger.Info().Msg("Cache miss, fetching from database")

		// If cache miss, retrieve users from the database.
		ctxTimeoutDB, cancelDB := usersUsecase.timeoutConfig.CreateDatabaseTimeout(ctx)
		defer cancelDB()

		usersFromDB, errDB := usersUsecase.usersRepository.GetAll(ctxTimeoutDB, request)
		if errDB != nil {
			usersUsecase.logger.Error().Msgf("Failed to fetch users from database: %v", errDB)
			return nil, usersUsecase.handleErrFromRepository(errDB, "Failed to fetch users from database: ")
		}

		// Set a timeout context for cache operations.
		ctxTimeoutSet, cancelSet := usersUsecase.timeoutConfig.CreateCacheTimeout(ctx)
		defer cancelSet() // Ensure context cancellation after function completes.

		// Cache the retrieved users data.
		errCacheSet := usersUsecase.cacheRepository.SetToCache(ctxTimeoutSet, key, usersFromDB)
		if errCacheSet != nil {
			// Log the error if caching fails but continue processing.
			usersUsecase.logger.Error().Msgf("Failed to set users to cache: %v", errCacheSet)
		}
		// Use data retrieved from the database.
		users = usersFromDB
	}

	// Set a timeout context for database count operation.
	ctxTimeoutDB, cancelDB := usersUsecase.timeoutConfig.CreateDatabaseTimeout(ctx)
	defer cancelDB()

	// Count the total number of users in the database.
	totalData, errCount := usersUsecase.usersRepository.Count(ctxTimeoutDB)
	if errCount != nil {
		usersUsecase.logger.Error().Msgf("Failed to count total number of users: %v", errCount)
		return nil, usersUsecase.handleErrFromRepository(errCount, "Failed to count total number of users: ")
	}

	// Calculate the total number of pages.
	totalPage := math.Ceil(float64(totalData) / float64(request.Size))

	// Prepare the response data with metadata and links.
	responseData := &response.LinksAble{
		Status: http.StatusOK,
		Code:   "STATUS_OK",
		Data:   mapper.EntitiesUserToResponses(users),
		Meta: map[string]any{
			"total_data": totalData,
			"total_page": totalPage,
			"size":       request.Size,
		},
		Links: map[string]any{
			"self": request.GetQueryParams(),
		},
	}

	usersUsecase.logger.Info().Msg("List method completed successfully")
	// Return the response data and any encountered errors.
	return responseData, nil
}

// Create handles the creation of a new user and associated operations.
func (usersUsecase UsersUsecase) Create(ctx context.Context, request *request.User) (*response.Standard, *response.StandardErrors) {
	usersUsecase.logger.Info().Msg("Create method called")

	var errHash error
	// Validate the incoming request data.
	if errValidate := usersUsecase.validator.Validate(request); errValidate != nil {
		usersUsecase.logger.Error().Msgf("Validation error: %v", errValidate)
		return nil, &response.StandardErrors{Errors: errValidate}
	}
	usersUsecase.logger.Info().Msg("Request validated successfully")

	// Set a timeout context for database existence check.
	ctxCount, cancelCount := usersUsecase.timeoutConfig.CreateDatabaseTimeout(ctx)
	defer cancelCount()

	// Check if a user with the same email already exists.
	exist, errExist := usersUsecase.usersRepository.ExistByKeyValue(ctxCount, map[string]any{"email": request.Email})
	if errExist != nil {
		usersUsecase.logger.Error().Msgf("Failed to count users in database: %v", errExist)
		return nil, usersUsecase.handleErrFromRepository(errExist, "Failed to count users in database")
	}
	// Return conflict error if the user already exists.
	if exist {
		usersUsecase.logger.Info().Msgf("User with email '%s' already exists", request.Email)
		return nil, &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.CONFLICT, "Users with email '"+request.Email+"' is exist")}}
	}

	// Generate a unique ID for the new user.
	id := ksuid.New().String()
	usersUsecase.logger.Info().Msgf("Generated new user ID: %s", id)

	// Hash the user's password.
	if request.Password, errHash = usersUsecase.hashing.Create(request.Password); errHash != nil {
		usersUsecase.logger.Error().Msgf("Failed to hash password: %v", errHash)
		return nil, &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.INTERNAL_SERVER_ERROR, "Internal Server Error : "+errHash.Error())}}
	}
	// Set a timeout context for database creation operation.
	ctxDB, cancel := usersUsecase.timeoutConfig.CreateDatabaseTimeout(ctx)
	defer cancel()

	// Save the new user to the database.
	if errDB := usersUsecase.usersRepository.Create(ctxDB, mapper.RequestUserToEntity(id, *request)); errDB != nil {
		usersUsecase.logger.Error().Msgf("Failed to save in database: %v", errDB)
		return nil, usersUsecase.handleErrFromRepository(errDB, "Failed to save in database")
	}

	// Invalidate related cache entries after database changes.
	if errCache := usersUsecase.handleDeleteCache(ctx, "users:*"); errCache != nil {
		usersUsecase.logger.Error().Msgf("Failed to invalidate cache: %v", errCache)
		return nil, errCache
	}

	// Retrieve the created user by ID for the response.
	users, standardErrors := usersUsecase.handleGetById(ctx, id)
	if standardErrors != nil {
		usersUsecase.logger.Error().Msgf("Failed to retrieve created user by ID: %v", standardErrors)
		return nil, standardErrors
	}

	// Return a successful response with the created user data.
	usersUsecase.logger.Info().Msg("User created successfully")
	return &response.Standard{
		Status: http.StatusCreated,
		Code:   "STATUS_CREATED",
		Data:   mapper.EntityUserToResponse(users),
	}, nil
}

// Update handles updating an existing user's information.
func (usersUsecase UsersUsecase) Update(ctx context.Context, request *request.User, id string) (*response.Standard, *response.StandardErrors) {
	usersUsecase.logger.Info().Msg("Update method called")

	var errHash error

	// Validate the user ID and request data.
	if errValidateVars := usersUsecase.validator.ValidateVars(id, "ksuid"); errValidateVars != nil {
		usersUsecase.logger.Error().Msgf("ID validation error: %v", errValidateVars)
		return nil, &response.StandardErrors{Errors: errValidateVars}
	}
	if errValidate := usersUsecase.validator.Validate(request); errValidate != nil {
		usersUsecase.logger.Error().Msgf("Validation error: %v", errValidate)
		return nil, &response.StandardErrors{Errors: errValidate}
	}
	usersUsecase.logger.Info().Msg("Request validated successfully")

	// Check if the user exists by ID.
	if errCount := usersUsecase.handleCountById(ctx, id); errCount != nil {
		usersUsecase.logger.Error().Msgf("User existence check failed: %v", errCount)
		return nil, errCount
	}

	// Hash the user's password if provided.
	if request.Password != "" {
		if request.Password, errHash = usersUsecase.hashing.Create(request.Password); errHash != nil {
			usersUsecase.logger.Error().Msgf("Failed to hash password: %v", errHash)
			return nil, &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.INTERNAL_SERVER_ERROR, "Internal Server Error : "+errHash.Error())}}
		}
	}

	// Set a timeout context for database update operation.
	ctxDB, cancel := usersUsecase.timeoutConfig.CreateDatabaseTimeout(ctx)
	defer cancel()

	// Update the user in the database.
	if errDB := usersUsecase.usersRepository.Update(ctxDB, mapper.RequestUserToEntity(id, *request), id); errDB != nil {
		usersUsecase.logger.Error().Msgf("Failed to update in database: %v", errDB)
		return nil, usersUsecase.handleErrFromRepository(errDB, "Failed to update in database")
	}

	// Invalidate related cache entries after database changes.
	if errCache := usersUsecase.handleDeleteCache(ctx, "users:*"); errCache != nil {
		usersUsecase.logger.Error().Msgf("Failed to invalidate cache: %v", errCache)
		return nil, errCache
	}

	// Retrieve the updated user by ID for the response.
	users, standardErrors := usersUsecase.handleGetById(ctx, id)
	if standardErrors != nil {
		usersUsecase.logger.Error().Msgf("Failed to retrieve updated user by ID: %v", standardErrors)
		return nil, standardErrors
	}

	// Return a successful response with the updated user data.
	usersUsecase.logger.Info().Msg("User updated successfully")
	return &response.Standard{
		Status: http.StatusOK,
		Code:   "STATUS_OK",
		Data:   mapper.EntityUserToResponse(users),
	}, nil
}

// Edit handles editing an existing user's information.
func (usersUsecase UsersUsecase) Edit(ctx context.Context, request *request.UserEdit, id string) (*response.Standard, *response.StandardErrors) {
	usersUsecase.logger.Info().Msg("Update method called")

	var errHash error

	// Validate the user ID and request data.
	if errValidateVars := usersUsecase.validator.ValidateVars(id, "ksuid"); errValidateVars != nil {
		usersUsecase.logger.Error().Msgf("ID validation error: %v", errValidateVars)
		return nil, &response.StandardErrors{Errors: errValidateVars}
	}
	if errValidate := usersUsecase.validator.Validate(request); errValidate != nil {
		usersUsecase.logger.Error().Msgf("Validation error: %v", errValidate)
		return nil, &response.StandardErrors{Errors: errValidate}
	}
	usersUsecase.logger.Info().Msg("Request validated successfully")

	// Check if the user exists by ID.
	if errCount := usersUsecase.handleCountById(ctx, id); errCount != nil {
		usersUsecase.logger.Error().Msgf("User existence check failed: %v", errCount)
		return nil, errCount
	}

	// Hash the user's password if provided.
	if request.Password != "" {
		if request.Password, errHash = usersUsecase.hashing.Create(request.Password); errHash != nil {
			usersUsecase.logger.Error().Msgf("Failed to hash password: %v", errHash)
			return nil, &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.INTERNAL_SERVER_ERROR, "Internal Server Error : "+errHash.Error())}}
		}
	}

	// Set a timeout context for database update operation.
	ctxDB, cancel := usersUsecase.timeoutConfig.CreateDatabaseTimeout(ctx)
	defer cancel()

	// Update the user in the database.
	if errDB := usersUsecase.usersRepository.Update(ctxDB, mapper.RequestUserEditToEntity(id, *request), id); errDB != nil {
		usersUsecase.logger.Error().Msgf("Failed to update in database: %v", errDB)
		return nil, usersUsecase.handleErrFromRepository(errDB, "Failed to update in database")
	}

	// Invalidate related cache entries after database changes.
	if errCache := usersUsecase.handleDeleteCache(ctx, "users:*"); errCache != nil {
		usersUsecase.logger.Error().Msgf("Failed to invalidate cache: %v", errCache)
		return nil, errCache
	}

	// Retrieve the updated user by ID for the response.
	users, standardErrors := usersUsecase.handleGetById(ctx, id)
	if standardErrors != nil {
		usersUsecase.logger.Error().Msgf("Failed to retrieve updated user by ID: %v", standardErrors)
		return nil, standardErrors
	}

	// Return a successful response with the updated user data.
	usersUsecase.logger.Info().Msg("User updated successfully")
	return &response.Standard{
		Status: http.StatusOK,
		Code:   "STATUS_OK",
		Data:   mapper.EntityUserToResponse(users),
	}, nil
}

// Delete removes a user by their ID.
func (usersUsecase UsersUsecase) Delete(ctx context.Context, id string) (*response.Standard, *response.StandardErrors) {
	usersUsecase.logger.Info().Msg("Delete method called")

	// Validate the user ID.
	if errValidate := usersUsecase.validator.ValidateVars(id, "ksuid"); errValidate != nil {
		usersUsecase.logger.Error().Msgf("ID validation error: %v", errValidate)
		return nil, &response.StandardErrors{Errors: errValidate}
	}
	usersUsecase.logger.Info().Msg("User ID validated successfully")

	// Check if the user exists by ID.
	if errCount := usersUsecase.handleCountById(ctx, id); errCount != nil {
		usersUsecase.logger.Error().Msgf("User existence check failed: %v", errCount)
		return nil, errCount
	}

	// Set a timeout context for database deletion operation.
	ctxDB, cancel := usersUsecase.timeoutConfig.CreateDatabaseTimeout(ctx)
	defer cancel()

	// Delete the user from the database.
	if errDB := usersUsecase.usersRepository.Delete(ctxDB, id); errDB != nil {
		usersUsecase.logger.Error().Msgf("Failed to delete from database: %v", errDB)
		return nil, usersUsecase.handleErrFromRepository(errDB, "Failed to delete from database")
	}

	// Invalidate related cache entries after database changes.
	if errCache := usersUsecase.handleDeleteCache(ctx, "users:*"); errCache != nil {
		usersUsecase.logger.Error().Msgf("Failed to invalidate cache: %v", errCache)
		return nil, errCache
	}

	// Return a successful response indicating the user was deleted.
	usersUsecase.logger.Info().Msg("User deleted successfully")
	return &response.Standard{
		Status: http.StatusOK,
		Code:   "STATUS_OK",
	}, nil
}

// Restore restore a user by their ID.
func (usersUsecase UsersUsecase) Restore(ctx context.Context, id string) (*response.Standard, *response.StandardErrors) {
	usersUsecase.logger.Info().Msg("Restore method called")

	// Validate the user ID.
	if errValidate := usersUsecase.validator.ValidateVars(id, "ksuid"); errValidate != nil {
		usersUsecase.logger.Error().Msgf("ID validation error: %v", errValidate)
		return nil, &response.StandardErrors{Errors: errValidate}
	}
	usersUsecase.logger.Info().Msg("User ID validated successfully")

	// Check if the user exists by ID
	// Set a timeout context for database count operation.
	ctxDBCount, cancelCount := usersUsecase.timeoutConfig.CreateDatabaseTimeout(ctx)
	defer cancelCount()

	count, err := usersUsecase.usersRepository.CountById(ctxDBCount, id)
	if err != nil {
		usersUsecase.logger.Error().Msgf("Failed to count user by ID: %v", err)
		return nil, usersUsecase.handleErrFromRepository(err, "Failed to count user by ID: ")
	}
	if count == 1 {
		usersUsecase.logger.Info().Msgf("User with ID '%s' is exist cannot restore", id)
		return nil, &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.NOT_FOUND, "User with ID '"+id+"' not found")}}
	}

	// Set a timeout context for database deletion operation.
	ctxDB, cancel := usersUsecase.timeoutConfig.CreateDatabaseTimeout(ctx)
	defer cancel()

	// Delete the user from the database.
	if errDB := usersUsecase.usersRepository.Restore(ctxDB, id); errDB != nil {
		usersUsecase.logger.Error().Msgf("Failed to delete from database: %v", errDB)
		return nil, usersUsecase.handleErrFromRepository(errDB, "Failed to delete from database")
	}

	// Invalidate related cache entries after database changes.
	if errCache := usersUsecase.handleDeleteCache(ctx, "users:*"); errCache != nil {
		usersUsecase.logger.Error().Msgf("Failed to invalidate cache: %v", errCache)
		return nil, errCache
	}

	// Return a successful response indicating the user was deleted.
	usersUsecase.logger.Info().Msg("User deleted successfully")
	return &response.Standard{
		Status: http.StatusOK,
		Code:   "STATUS_OK",
	}, nil
}

// Get used for get uses by id
func (usersUsecase UsersUsecase) Get(ctx context.Context, id string) (*response.Standard, *response.StandardErrors) {
	usersUsecase.logger.Info().Msg("Get method called")

	// Validate the ID format
	if errValidateVars := usersUsecase.validator.ValidateVars(id, "ksuid"); errValidateVars != nil {
		usersUsecase.logger.Error().Msgf("ID validation error: %v", errValidateVars)
		return nil, &response.StandardErrors{Errors: errValidateVars}
	}

	// Check if the user exists by ID
	if errCount := usersUsecase.handleCountById(ctx, id); errCount != nil {
		usersUsecase.logger.Error().Msgf("User existence check failed: %v", errCount)
		return nil, errCount
	}

	// Retrieve user details by ID
	users, standardErrors := usersUsecase.handleGetById(ctx, id)
	if standardErrors != nil {
		usersUsecase.logger.Error().Msgf("Failed to retrieve get user by ID: %v", standardErrors)
		return nil, standardErrors
	}

	// Return user details with HTTP 200 OK status
	return &response.Standard{
		Status: http.StatusOK,
		Code:   "STATUS_OK",
		Data:   mapper.EntityUserToResponse(users),
	}, nil
}

// handleCountById checks if a user with the specified ID exists in the database.
func (usersUsecase UsersUsecase) handleCountById(ctx context.Context, id string) *response.StandardErrors {
	usersUsecase.logger.Info().Msg("handleCountById method called")

	// Set a timeout context for database count operation.
	ctxDB, cancel := usersUsecase.timeoutConfig.CreateDatabaseTimeout(ctx)
	defer cancel()

	// Check if the user exists by ID.
	count, err := usersUsecase.usersRepository.CountById(ctxDB, id)
	if err != nil {
		usersUsecase.logger.Error().Msgf("Failed to count user by ID: %v", err)
		return usersUsecase.handleErrFromRepository(err, "Failed to count user by ID: ")
	}
	if count != 1 {
		usersUsecase.logger.Info().Msgf("User with ID '%s' does not exist", id)
		return &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.NOT_FOUND, "User with ID '"+id+"' not found")}}
	}

	usersUsecase.logger.Info().Msgf("User with ID '%s' exists", id)
	return nil
}

// handleGetById retrieves a user by their ID from the database.
func (usersUsecase UsersUsecase) handleGetById(ctx context.Context, id string) (*entity.Users, *response.StandardErrors) {
	usersUsecase.logger.Info().Msg("handleGetById method called")

	// Set a timeout context for database retrieval operation.
	ctxDB, cancel := usersUsecase.timeoutConfig.CreateDatabaseTimeout(ctx)
	defer cancel()

	// Retrieve the user from the database.
	user := entity.Users{}
	err := usersUsecase.usersRepository.GetById(ctxDB, &user, id)
	if err != nil {
		usersUsecase.logger.Error().Msgf("Failed to fetch user by ID from database: %v", err)
		return nil, usersUsecase.handleErrFromRepository(err, "Failed to fetch user by ID from database: ")
	}

	usersUsecase.logger.Info().Msgf("User with ID '%s' retrieved successfully", id)
	return &user, nil
}

// handleErrFromRepository handles errors from the repository, including context.DeadlineExceeded, and logs them.
func (usersUsecase UsersUsecase) handleErrFromRepository(err error, message string) *response.StandardErrors {
	if errors.Is(err, context.DeadlineExceeded) {
		usersUsecase.logger.Error().Msgf("%s: operation timed out: %v", message, err)
		return &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.REQUEST_TIMEOUT, "Request timed out: "+err.Error())}}
	}

	usersUsecase.logger.Error().Msgf("%s: %v", message, err)
	return &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.INTERNAL_SERVER_ERROR, message+err.Error())}}
}

// handleDeleteCache invalidates cache entries related to users.
func (usersUsecase UsersUsecase) handleDeleteCache(ctx context.Context, pattern string) *response.StandardErrors {
	usersUsecase.logger.Info().Msgf("handleDeleteCache method called with pattern: %s", pattern)
	ctxTimeout, cancel := usersUsecase.timeoutConfig.CreateCacheTimeout(ctx)
	defer cancel()
	// Invalidate cache entries based on the pattern.
	if err := usersUsecase.cacheRepository.DeleteToCacheByRegexKey(ctxTimeout, pattern); err != nil {
		usersUsecase.logger.Error().Msgf("Failed to invalidate cache with pattern '%s': %v", pattern, err)
		return usersUsecase.handleErrFromRepository(err, "Failed to invalidate cache: "+err.Error())
	}

	usersUsecase.logger.Info().Msgf("Cache invalidated successfully for pattern: %s", pattern)
	return nil
}
