package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/phuslu/log"
	"github.com/segmentio/ksuid"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/hash"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/timeout"
	tokenconfig "github.com/tirtahakimpambudhi/restful_api/internal/configs/token"
	"github.com/tirtahakimpambudhi/restful_api/internal/entity"
	errorshandler "github.com/tirtahakimpambudhi/restful_api/internal/errors"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/request"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/response"
	"github.com/tirtahakimpambudhi/restful_api/internal/repository"
	"github.com/tirtahakimpambudhi/restful_api/internal/validation"
	"net/http"
	"time"
)

// AuthUsecase handles the authentication logic.
type AuthUsecase struct {
	usersRepository repository.UsersRepository // Repository to access user data.
	timeoutConfig   *timeout.Config            // Configuration for handling timeouts.
	validator       *validation.Validator      // Validator for request validation.
	hashing         *hash.Argon2               // Password hashing utility.
	token           *tokenconfig.JWTToken      // Token generation and verification utility.
	secretKey       *tokenconfig.SecretKey     // Secret key for Token secret
	logger          *log.Logger                // Logger for logging messages.
	enforcer        *casbin.Enforcer
}

// Login used for users login logic.
func (a AuthUsecase) Login(ctx context.Context, req *request.Auth) (*response.Standard, string, *response.StandardErrors) {
	a.logger.Info().Msg("Login method called") // Log the method call.

	// Validate the incoming request data.
	if errValidate := a.validator.Validate(req); errValidate != nil {
		a.logger.Error().Msgf("Validation error: %v", errValidate)    // Log validation error.
		return nil, "", &response.StandardErrors{Errors: errValidate} // Return validation errors.
	}
	a.logger.Info().Msg("Request validated successfully") // Log successful validation.

	// Set a timeout context for database existence check.
	ctxCount, cancelCount := a.timeoutConfig.CreateDatabaseTimeout(ctx)
	defer cancelCount() // Ensure the context is canceled.

	// Check if a user with the same email already exists.
	exist, errExist := a.usersRepository.ExistByKeyValue(ctxCount, map[string]any{"email": req.Email})
	if errExist != nil {
		a.logger.Error().Msgf("Failed to count users in database: %v", errExist)                 // Log error.
		return nil, "", a.handleErrFromRepository(errExist, "Failed to count users in database") // Handle repository error.
	}

	// Return conflict error if the user not exists.
	if !exist {
		a.logger.Info().Msgf("User with email '%s' not exists", req.Email)                                                                                                  // Log user not exists.
		return nil, "", &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.NOT_FOUND, "Users with email '"+req.Email+"' not exists")}} // Return not found error.
	}

	// Set a timeout context for database get users check.
	ctxGet, cancelGet := a.timeoutConfig.CreateDatabaseTimeout(ctx)
	defer cancelGet() // Ensure the context is canceled.

	// Retrieve the user by email.
	users, errGet := a.handleGetByEmail(ctxGet, req.Email)
	if errGet != nil {
		a.logger.Error().Msgf("Failed to get by email users in database: %v", errGet) // Log error.
		return nil, "", errGet                                                        // Return the error.
	}
	// Match the request password with users database.
	match, errMatch := a.hashing.Match(req.Password, users.Password)
	if errMatch != nil {
		a.logger.Error().Msgf("Failed to match by hashing users in database: %v", errMatch)                                                                                                                         // Log hashing error.
		return nil, "", &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.INTERNAL_SERVER_ERROR, fmt.Sprintf("Failed to match by hashing users in database: %v", errMatch))}} // Return hashing error.
	}

	// Handle password mismatch.
	if !match {
		a.logger.Error().Msg("Password users not match")                                                                                                  // Log password mismatch.
		return nil, "", &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.UNAUTHORIZE, "email or password wrong")}} // Return unauthorized error.
	}

	// Parse user ID.
	userId, errParse := ksuid.Parse(users.ID)
	if errParse != nil {
		a.logger.Error().Msgf("Failed to parse user ID: %v", errParse)                                                                                                         // Log parsing error.
		return nil, "", &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.INTERNAL_SERVER_ERROR, "Error Parse ID : "+errParse.Error())}} // Return parsing error.
	}

	// Set token expiration time.
	expiredAt := time.Now().Add(5 * time.Minute)
	a.logger.Info().Msgf("Successfully authenticated user with email '%s'", req.Email) // Log successful authentication.

	// Create access token.
	payload := tokenconfig.NewTokenPayloadBuilder().WithEmail(users.Email).WithUserID(userId).WithExpiration(expiredAt).Build()
	accessToken, errToken := a.handleCreateToken(a.secretKey.AccessToken, *payload)
	if errToken != nil {
		a.logger.Error().Msgf("Failed to create access token: %v", errToken) // Log token creation error.
		return nil, "", errToken                                             // Return token creation error.
	}

	// Set refresh token expiration time.
	expiredRefreshToken := time.Now().Add(7 * 24 * time.Hour)
	payloadRefreshToken := tokenconfig.NewTokenPayloadBuilder().WithEmail(users.Email).WithUserID(userId).WithExpiration(expiredRefreshToken).Build()
	// Create refresh token.
	refreshToken, errRefreshToken := a.handleCreateToken(a.secretKey.RefreshToken, *payloadRefreshToken)
	if errRefreshToken != nil {
		a.logger.Error().Msgf("Failed to create access token: %v", errRefreshToken) // Log refresh token creation error.
		return nil, "", errToken                                                    // Return refresh token creation error.
	}

	// Return the generated tokens.
	return &response.Standard{
		Status: http.StatusOK,
		Code:   "STATUS_OK",
		Data: &response.Token{
			AccessToken: accessToken,
			ExpiredAt:   expiredAt.UnixMilli(),
		},
	}, refreshToken, nil
}

// Logout handles user logout logic.
func (a AuthUsecase) Logout(ctx context.Context, token string) (*response.Standard, *response.StandardErrors) {
	a.logger.Info().Msg("Logout method called") // Log the method call.

	// Parse the token.
	_, standardErrors := a.handleParseToken(a.secretKey.RefreshToken, token)
	if standardErrors != nil {
		a.logger.Error().Msgf("Failed to parse token: %v", standardErrors) // Log token parsing error.
		return nil, standardErrors                                         // Return the token parsing error.
	}

	// Return successful logout response.
	return &response.Standard{
		Status: http.StatusOK,
		Code:   "STATUS_OK",
		Data:   nil,
	}, nil
}

// RefreshToken handles the logic for refreshing a user's token.
func (a AuthUsecase) RefreshToken(ctx context.Context, token string) (*response.Standard, *response.StandardErrors) {
	a.logger.Info().Msg("RefreshToken method called") // Log the method call.

	// Parse the token.
	payload, standardErrors := a.handleParseToken(a.secretKey.RefreshToken, token)
	if standardErrors != nil {
		a.logger.Error().Msgf("Failed to parse token: %v", standardErrors) // Log token parsing error.
		return nil, standardErrors                                         // Return token parsing error.
	}

	// Set the new token expiration time.
	expiredAt := time.Now().Add(5 * time.Minute)

	// Create new access token.
	accessToken, errToken := a.handleCreateToken(a.secretKey.AccessToken, *tokenconfig.NewTokenPayloadBuilder().WithEmail(payload.Email).WithUserID(payload.ID).WithExpiration(expiredAt).Build())
	if errToken != nil {
		a.logger.Error().Msgf("Failed to create access token: %v", errToken) // Log token creation error.
		return nil, errToken                                                 // Return token creation error.
	}

	// Return the new access token.
	return &response.Standard{
		Status: http.StatusOK,
		Code:   "STATUS_OK",
		Data: &response.Token{
			AccessToken: accessToken,
			ExpiredAt:   expiredAt.UnixMilli(),
		},
	}, nil
}

// ForgotPassword handles the logic for the forgot password functionality.
func (a AuthUsecase) ForgotPassword(ctx context.Context, req *request.ForgotPassword) (*response.Standard, *response.StandardErrors) {
	a.logger.Info().Msg("ForgotPassword method called") // Log the method call.

	// Validate the incoming request data
	if errValidate := a.validator.Validate(req); errValidate != nil {
		a.logger.Error().Msgf("Validation error: %v", errValidate)
		return nil, &response.StandardErrors{Errors: errValidate}
	}

	// Set a timeout context for database existence check.
	ctxCount, cancelCount := a.timeoutConfig.CreateDatabaseTimeout(ctx)
	defer cancelCount() // Ensure the context is canceled.

	// Check if a user with the same email already exists.
	exist, errExist := a.usersRepository.ExistByKeyValue(ctxCount, map[string]any{"email": req.Email})
	if errExist != nil {
		a.logger.Error().Msgf("Failed to count users in database: %v", errExist)             // Log counting error.
		return nil, a.handleErrFromRepository(errExist, "Failed to count users in database") // Handle repository error.
	}

	// Return conflict error if the user not exists.
	if !exist {
		a.logger.Info().Msgf("User with email '%s' not exists", req.Email)                                                                                             // Log user not exists.
		return nil, &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.NOT_FOUND, "Users with email '"+req.Email+"' not exist")}} // Return not found error.
	}

	a.logger.Info().Msgf("Successfully forgot password for user '%s' with email", req.Email) // Log successful forgot password.
	// TODO : Publish in Downstream. and response token

	// Return success response.
	return &response.Standard{
		Status: http.StatusOK,
		Code:   "STATUS_OK",
		Data:   nil,
	}, nil
}

// ResetPassword handles the reset password process by validating the request, parsing the token,
// and updating the user's password in the database.
func (a AuthUsecase) ResetPassword(ctx context.Context, payload *tokenconfig.Payload, req *request.ResetPassword) (*response.Standard, *response.StandardErrors) {
	var errHash error

	// Log the method call
	a.logger.Info().Msg("ResetPassword method called")

	// Validate the incoming request data
	if errValidate := a.validator.Validate(req); errValidate != nil {
		a.logger.Error().Msgf("Validation error: %v", errValidate)
		return nil, &response.StandardErrors{Errors: errValidate}
	}

	// Set a timeout context for fetching the user by email
	ctxGet, cancelGet := a.timeoutConfig.CreateDatabaseTimeout(ctx)
	defer cancelGet()

	// Fetch the user by email from the database
	users, errGet := a.handleGetByEmail(ctxGet, payload.Email)
	if errGet != nil {
		a.logger.Error().Msgf("Failed to get user by email in database: %v", errGet)
		return nil, errGet
	}

	// Set a timeout context for updating the user's password
	ctxUpdate, cancelUpdate := a.timeoutConfig.CreateDatabaseTimeout(ctx)
	defer cancelUpdate()

	// Hash the new password
	if users.Password, errHash = a.hashing.Create(users.Password); errHash != nil {
		a.logger.Error().Msgf("Failed to hash password: %v", errHash)
		return nil, &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.INTERNAL_SERVER_ERROR, "Internal Server Error: "+errHash.Error())}}
	}

	// Update the user in the database
	if errDB := a.usersRepository.Update(ctxUpdate, users, users.ID); errDB != nil {
		a.logger.Error().Msgf("Failed to update user in database: %v", errDB)
		return nil, a.handleErrFromRepository(errDB, "Failed to update user in database")
	}

	// Return success response
	return &response.Standard{
		Status: http.StatusOK,
		Code:   "STATUS_OK",
		Data:   nil,
	}, nil
}

// handleGetByEmail retrieves a user by their email from the database.
func (a AuthUsecase) handleGetByEmail(ctx context.Context, email string) (*entity.Users, *response.StandardErrors) {
	a.logger.Info().Msg("handleGetByEmail method called")

	// Set a timeout context for the database retrieval operation
	ctxDB, cancel := a.timeoutConfig.CreateDatabaseTimeout(ctx)
	defer cancel()

	// Retrieve the user from the database
	user := entity.Users{}
	err := a.usersRepository.GetByEmail(ctxDB, &user, email)
	if err != nil {
		a.logger.Error().Msgf("Failed to fetch user by email from database: %v", err)
		return nil, a.handleErrFromRepository(err, "Failed to fetch user by email from database")
	}

	// Log success message
	a.logger.Info().Msgf("User with Email '%s' retrieved successfully", email)
	return &user, nil
}

// handleErrFromRepository handles errors from the repository, including context.DeadlineExceeded, and logs them.
func (a AuthUsecase) handleErrFromRepository(err error, message string) *response.StandardErrors {
	if errors.Is(err, context.DeadlineExceeded) {
		a.logger.Error().Msgf("%s: operation timed out: %v", message, err)
		return &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.REQUEST_TIMEOUT, "Request timed out: "+err.Error())}}
	}

	a.logger.Error().Msgf("%s: %v", message, err)
	return &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.INTERNAL_SERVER_ERROR, message+err.Error())}}
}

// handleCreateToken generates a new token based on the provided payload.
func (a AuthUsecase) handleCreateToken(secretKey string, payload tokenconfig.Payload) (string, *response.StandardErrors) {
	a.logger.Info().Msg("handleCreateToken method called")

	// Create a new token
	createToken, err := a.token.CreateToken(secretKey, &payload)
	if err != nil {
		return "", &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.INTERNAL_SERVER_ERROR, "Error Create Token: "+err.Error())}}
	}

	// Return the generated token
	return createToken, nil
}

// handleParseToken parses and verifies a token and returns its payload.
func (a AuthUsecase) handleParseToken(secretKey string, token string) (*tokenconfig.Payload, *response.StandardErrors) {
	a.logger.Info().Msg("handleParseToken method called")

	// Parse the token
	payload, errParse := a.token.VerifyToken(secretKey, token)
	if errParse != nil {
		a.logger.Error().Msgf("Failed to parse token: %v", errParse)

		var tokenErr *tokenconfig.TokenError
		if errors.As(errParse, &tokenErr) {
			err := new(response.StandardErrors)
			typeErr := tokenErr.TypeError()

			// Handle different types of token errors
			switch {
			case errors.Is(typeErr, tokenconfig.ErrInvalidToken):
				err.Errors = []*response.Error{
					errorshandler.NewError(errorshandler.BAD_REQUEST, "Error Invalid Token '"+token+"'"),
				}
			case errors.Is(typeErr, tokenconfig.ErrExpiredToken):
				err.Errors = []*response.Error{
					errorshandler.NewError(errorshandler.FORBIDEN, "Error Expired Token '"+token+"'"),
				}
			case errors.Is(typeErr, tokenconfig.ErrTokenMalformed):
				err.Errors = []*response.Error{
					errorshandler.NewError(errorshandler.BAD_REQUEST, "Error Token Malformed '"+token+"'"),
				}
			case errors.Is(typeErr, tokenconfig.ErrTokenSignatureInvalid):
				err.Errors = []*response.Error{
					errorshandler.NewError(errorshandler.BAD_REQUEST, "Error Invalid Token Signature '"+token+"'"),
				}
			case errors.Is(typeErr, tokenconfig.ErrFailedParseClaims):
				err.Errors = []*response.Error{
					errorshandler.NewError(errorshandler.UNPROCESS_ENITITY, "Error Parsing Claims '"+token+"'"),
				}
			case errors.Is(typeErr, tokenconfig.ErrInvalidKey):
				err.Errors = []*response.Error{
					errorshandler.NewError(errorshandler.UNPROCESS_ENITITY, "Error Invalid Secret Key '"+token+"'"),
				}
			default:
				err.Errors = []*response.Error{
					errorshandler.NewError(errorshandler.INTERNAL_SERVER_ERROR, "Server Error"),
				}
			}
			return nil, err
		}
	}
	// Return the parsed payload
	return payload, nil
}

// UpdateRole used for update or insert role in grouping
func (a AuthUsecase) UpsertRole(ctx context.Context, req *request.UpdateRole) (*response.Standard, *response.StandardErrors) {
	a.logger.Info().Msg("UpsertRole method called")

	// Validate the incoming request data
	if errValidate := a.validator.Validate(req); errValidate != nil {
		a.logger.Error().Msgf("Validation error: %v", errValidate)
		return nil, &response.StandardErrors{Errors: errValidate}
	}

	// Get Role for user by email
	roles, errRole := a.enforcer.GetRolesForUser(req.Email)
	if errRole != nil {
		return nil, &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.INTERNAL_SERVER_ERROR, "Error Internal Server: "+errRole.Error())}}
	}
	// check the role if greater than 0 remove it
	if len(roles) > 0 {
		for _, role := range roles {
			removed, errRemoved := a.enforcer.DeleteRoleForUser(req.Email, role)
			if errRemoved != nil {
				return nil, &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.INTERNAL_SERVER_ERROR, "Error Internal Server: "+errRemoved.Error())}}
			}
			if removed {
				a.logger.Info().Msgf("Removed existing role %s from user %s", role, req.Email)
			}
		}
	}

	// Add role in user
	added, errAdd := a.enforcer.AddGroupingPolicy(req.Email, req.RoleName)
	if added {
		a.logger.Info().Msgf("Role %s added to user %s", req.RoleName, req.Email)
	}

	if errAdd != nil {
		return nil, &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.INTERNAL_SERVER_ERROR, "Error Internal Server: "+errAdd.Error())}}
	}
	return &response.Standard{
		Status: http.StatusOK,
		Code:   "STATUS_OK",
		Data:   map[string]any{"message": "Successfully Upsert Role"},
	}, nil
}
