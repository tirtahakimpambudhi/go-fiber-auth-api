// Package http provides HTTP handlers for user operations
package http

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/phuslu/log"
	errorshandler "github.com/tirtahakimpambudhi/restful_api/internal/errors"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/request"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/response"
	"github.com/tirtahakimpambudhi/restful_api/internal/usecase"
	reflecthelper "github.com/tirtahakimpambudhi/restful_api/pkg/helper/reflect"
)

// UsersController handles requests related to user operations
type UsersController struct {
	usecases *usecase.UsersUsecase
	logger   *log.Logger
}

// NewUsersController creates a new UsersController
func NewUsersController(usecases *usecase.UsersUsecase, logger *log.Logger) *UsersController {
	// Initialize UsersController with provided usecases
	logger.Info().Msg("UsersController initialized")
	return &UsersController{usecases: usecases, logger: logger}
}

// Index retrieves a list of users
func (controller UsersController) Index(ctx *fiber.Ctx) error {
	// Log the start of the Index method
	controller.logger.Info().Msg("Index method called")

	// Create a new Page request
	req := new(request.Page)
	controller.logger.Debug().Msg("Created new Page request")

	// Parse query parameters into the request struct
	if errParse := ctx.QueryParser(req); errParse != nil {
		// Log the error during parsing
		controller.logger.Error().Err(errParse).Msg("Failed to parse query parameters")
		// Return a bad request error if parsing fails
		return &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.BAD_REQUEST, fmt.Sprintf("BAD REQUEST : %s \nREQUEST BODY \n%s", errParse.Error(), reflecthelper.KeyValueToString(*req)))}}
	}
	fmt.Println(req)
	// Retrieve a list of users from the usecase
	controller.logger.Info().Msg("Calling usecase List method")
	res, errors := controller.usecases.List(ctx.Context(), req)
	if errors != nil {
		// Log the error during user retrieval
		controller.logger.Error().Err(errors).Msg("Failed to retrieve user list")
		// Return any errors encountered during retrieval
		return errors
	}

	// If there are links, update them to include the full URL
	if len(res.Links) > 0 {
		baseURL := ctx.BaseURL()
		for key, value := range res.Links {
			// Build the full URL using baseURL, path, and query parameters
			fullURL := fmt.Sprintf("%s%s", baseURL, value)
			res.Links[key] = fullURL
			controller.logger.Debug().Msgf("Updated link %s to %s", key, fullURL)
		}
	}

	// Log the successful retrieval of users
	controller.logger.Info().Msg("Successfully retrieved users")

	// Set the response status code
	ctx.Status(res.Status)

	// Return the response as JSON
	return ctx.JSON(res)
}

// Show retrieves a single user by ID
func (controller UsersController) Show(ctx *fiber.Ctx) error {
	// Log the start of the Show method
	controller.logger.Info().Msg("Show method called")

	// Retrieve the user by ID from the usecase
	controller.logger.Info().Msg("Calling usecase Get method")
	res, errors := controller.usecases.Get(ctx.Context(), utils.CopyString(ctx.Params("id")))
	if errors != nil {
		// Log the error during user retrieval
		controller.logger.Error().Err(errors).Msg("Failed to retrieve user")
		// Return any errors encountered during retrieval
		return errors
	}

	// Log the successful retrieval of the user
	controller.logger.Info().Msgf("Successfully retrieved user with ID: %s", ctx.Params("id"))

	// Set the response status code
	ctx.Status(res.Status)

	// Return the response as JSON
	return ctx.JSON(res)
}

// Store creates a new user
func (controller UsersController) Store(ctx *fiber.Ctx) error {
	// Log the start of the Store method
	controller.logger.Info().Msg("Store method called")

	// Create a new User request
	req := new(request.User)
	controller.logger.Debug().Msg("Created new User request")

	// Parse the request body into the User struct
	if errParse := ctx.BodyParser(req); errParse != nil {
		// Log the error during body parsing
		controller.logger.Error().Err(errParse).Msg("Failed to parse request body")
		// Return a bad request error if parsing fails
		return &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.BAD_REQUEST, fmt.Sprintf("BAD REQUEST : %s \nREQUEST BODY \n%s", errParse.Error(), reflecthelper.KeyValueToString(*req)))}}
	}

	// Log the received request data
	controller.logger.Debug().Msgf("Received request data: %+v", req)

	// Create a new user using the usecase
	controller.logger.Info().Msg("Calling usecase Create method")
	res, errors := controller.usecases.Create(ctx.Context(), req)
	if errors != nil {
		// Log the error during user creation
		controller.logger.Error().Err(errors).Msg("Failed to create user")
		// Return any errors encountered during creation
		return errors
	}

	// Log the successful creation of the user
	controller.logger.Info().Msg("Successfully created new user")

	// Set the response status code
	ctx.Status(res.Status)

	// Return the response as JSON
	return ctx.JSON(res)
}

// Update modifies an existing user by ID
func (controller UsersController) Update(ctx *fiber.Ctx) error {
	// Log the start of the Update method
	controller.logger.Info().Msg("Update method called")

	// Create a new User request
	req := new(request.User)
	controller.logger.Debug().Msg("Created new User request")

	// Parse the request body into the User struct
	if errParse := ctx.BodyParser(req); errParse != nil {
		// Log the error during body parsing
		controller.logger.Error().Err(errParse).Msg("Failed to parse request body")
		// Return a bad request error if parsing fails
		return &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.BAD_REQUEST, fmt.Sprintf("BAD REQUEST : %s \nREQUEST BODY \n%s", errParse.Error(), reflecthelper.KeyValueToString(*req)))}}
	}

	// Log the received request data
	controller.logger.Debug().Msgf("Received request data: %+v", req)

	// Update the user using the usecase
	controller.logger.Info().Msg("Calling usecase Update method")
	res, errors := controller.usecases.Update(ctx.Context(), req, utils.CopyString(ctx.Params("id")))
	if errors != nil {
		// Log the error during user update
		controller.logger.Error().Err(errors).Msg("Failed to update user")
		// Return any errors encountered during update
		return errors
	}

	// Log the successful update of the user
	controller.logger.Info().Msgf("Successfully updated user with ID: %s", ctx.Params("id"))

	// Set the response status code
	ctx.Status(res.Status)

	// Return the response as JSON
	return ctx.JSON(res)
}

// Edit modifies specific fields of an existing user by ID
func (controller UsersController) Edit(ctx *fiber.Ctx) error {
	// Log the start of the Edit method
	controller.logger.Info().Msg("Edit method called")

	// Create a new UserEdit request
	req := new(request.UserEdit)
	controller.logger.Debug().Msg("Created new UserEdit request")

	// Parse the request body into the UserEdit struct
	if errParse := ctx.BodyParser(req); errParse != nil {
		// Log the error during body parsing
		controller.logger.Error().Err(errParse).Msg("Failed to parse request body")
		// Return a bad request error if parsing fails
		return &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.BAD_REQUEST, fmt.Sprintf("BAD REQUEST : %s \nREQUEST BODY \n%s", errParse.Error(), reflecthelper.KeyValueToString(*req)))}}
	}

	// Log the received request data
	controller.logger.Debug().Msgf("Received request data: %+v", req)

	// Edit the user using the usecase
	controller.logger.Info().Msg("Calling usecase Edit method")
	res, errors := controller.usecases.Edit(ctx.Context(), req, utils.CopyString(ctx.Params("id")))
	if errors != nil {
		// Log the error during user edit
		controller.logger.Error().Err(errors).Msg("Failed to edit user")
		// Return any errors encountered during editing
		return errors
	}

	// Log the successful editing of the user
	controller.logger.Info().Msgf("Successfully edited user with ID: %s", ctx.Params("id"))

	// Set the response status code
	ctx.Status(res.Status)

	// Return the response as JSON
	return ctx.JSON(res)
}

// Destroy deletes a user by ID
func (controller UsersController) Destroy(ctx *fiber.Ctx) error {
	// Log the start of the Destroy method
	controller.logger.Info().Msg("Destroy method called")

	// Delete the user using the usecase
	controller.logger.Info().Msg("Calling usecase Delete method")
	res, errors := controller.usecases.Delete(ctx.Context(), utils.CopyString(ctx.Params("id")))
	if errors != nil {
		// Log the error during user deletion
		controller.logger.Error().Err(errors).Msg("Failed to delete user")
		// Return any errors encountered during deletion
		return errors
	}

	// Log the successful deletion of the user
	controller.logger.Info().Msgf("Successfully deleted user with ID: %s", ctx.Params("id"))

	// Set the response status code
	ctx.Status(res.Status)

	// Return the response as JSON
	return ctx.JSON(res)
}

// Restore restore a user by ID was a delete
func (controller UsersController) Restore(ctx *fiber.Ctx) error {
	// Log the start of the Restore method
	controller.logger.Info().Msg("Restore method called")

	// Delete the user using the usecase
	controller.logger.Info().Msg("Calling usecase Restore Usecases method")
	res, errors := controller.usecases.Restore(ctx.Context(), utils.CopyString(ctx.Params("id")))
	if errors != nil {
		// Log the error during user restoration
		controller.logger.Error().Err(errors).Msg("Failed to restore user")
		// Return any errors encountered during restoration
		return errors
	}

	// Log the successful restore of the user
	controller.logger.Info().Msgf("Successfully restore user with ID: %s", ctx.Params("id"))

	// Set the response status code
	ctx.Status(res.Status)

	// Return the response as JSON
	return ctx.JSON(res)
}
