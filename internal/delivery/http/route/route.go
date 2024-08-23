package route

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
	fiberlog "github.com/phuslu/log/fiber"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/bootstrap"
	loggerconfig "github.com/tirtahakimpambudhi/restful_api/internal/configs/logger"
	tokenconfig "github.com/tirtahakimpambudhi/restful_api/internal/configs/token"
	"github.com/tirtahakimpambudhi/restful_api/internal/delivery/http"
	"github.com/tirtahakimpambudhi/restful_api/internal/delivery/http/middleware"
	errorshandler "github.com/tirtahakimpambudhi/restful_api/internal/errors"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/response"
)

// Route struct holds controllers and the application configuration
type Route struct {
	UsersController  *http.UsersController
	AuthController   *http.AuthController
	Logger           *loggerconfig.Logger
	CasbinMiddleware *casbin.Enforcer
	Token            *tokenconfig.JWTToken
	SecretKey        *tokenconfig.SecretKey
}

// NewRoute initializes and returns a new Route instance
func NewRoute(app *bootstrap.App, usersController *http.UsersController, authController *http.AuthController) (*Route, error) {
	app.Logger.App.Info().Msg("NewRoute Call Function")

	// Create a new Route instance with the controllers and app configuration
	routes := &Route{UsersController: usersController, AuthController: authController, SecretKey: app.Secret, Logger: app.Logger, Token: app.Token, CasbinMiddleware: app.CasbinEnforcer}

	// Return the initialized Route instance
	return routes, nil
}

func (r *Route) Init(app *fiber.App) error {

	// Setup middleware for the Fiber application
	if err := middleware.Setup(app); err != nil {
		r.Logger.App.Error().Err(err)
		return err
	}
	group := app.Group("/api/v1")
	group.Get("/monitor", middleware.Monitor())
	r.public(group)
	r.protected(group)
	group.Use(fiberlog.New(r.Logger.Access, nil)) // WARNING : ALWAYS BEFORE NOT FOUND HANDLER
	app.Use(func(ctx *fiber.Ctx) error {
		method := ctx.Method() // Get Method from request client
		path := ctx.Path()     // Get Path from request client

		// Set status to 404 Not Found
		ctx.Status(fiber.StatusNotFound)

		// Create message error with specific method dan path
		errorMessage := fmt.Sprintf("NOT FOUND: Method %s with Path %s", method, path)

		// Send error with detail error
		return ctx.JSON(&response.StandardErrors{
			Errors: []*response.Error{
				errorshandler.NewError(errorshandler.NOT_FOUND, errorMessage),
			},
		})
	})

	return nil
}

// Public sets up the public routes
func (r *Route) public(group fiber.Router) {
	r.Logger.App.Info().Msg("Prepare the public routes")
	// Define a route for creating a new user
	group.Post("/users", r.UsersController.Store)

	// Define routes for authentication
	authRoute := group.Group("/auth")
	authRoute.Post("/login", r.AuthController.Login)
	authRoute.Delete("/logout", r.AuthController.Logout)
	authRoute.Get("/refresh-token", r.AuthController.RefreshToken)
}

// Protected sets up the protected routes with middleware
func (r *Route) protected(group fiber.Router) {
	r.Logger.App.Info().Msg("Prepare the protected routes")
	// Define a route for resetting the password, protected by a middleware
	group.Post("/auth/reset-password", middleware.NewAuthenticationToken(r.Token, r.SecretKey.ForgotPasswordToken), r.AuthController.ResetPassword)
	group.Patch("/auth/role", middleware.NewAuthenticationToken(r.Token, r.SecretKey.AccessToken), middleware.NewAuthorization(r.CasbinMiddleware, "admin"), r.AuthController.UpsertRole)
	// Define a group of routes protected by access token authentication
	usersProtectedRoute := group.Group("/users", middleware.NewAuthenticationToken(r.Token, r.SecretKey.AccessToken))

	// Define a route for getting all users with required permissions
	usersProtectedRoute.Get("", middleware.NewAuthorizationById(r.CasbinMiddleware, "users:read"), r.UsersController.Index)

	// Define a route for getting a specific user by ID, protected by custom and Casbin middleware
	usersProtectedRoute.Get("/:id", middleware.NewAuthorizationById(r.CasbinMiddleware, "users:read"), r.UsersController.Show)

	// Define a route for updating a user by ID, protected by custom middleware
	usersProtectedRoute.Put("/:id", middleware.NewAuthorizationById(r.CasbinMiddleware, "users:update"), r.UsersController.Update)

	// Define a route for editing a user by ID, protected by custom middleware
	usersProtectedRoute.Patch("/:id", middleware.NewAuthorizationById(r.CasbinMiddleware, "users:edit"), r.UsersController.Edit)

	// Define a route for deleting a user by ID, protected by Casbin middleware
	usersProtectedRoute.Delete("/:id", middleware.NewAuthorization(r.CasbinMiddleware, "admin"), r.UsersController.Destroy)

	// Define a route for deleting a user by ID, protected by Casbin middleware
	usersProtectedRoute.Post("/:id", middleware.NewAuthorization(r.CasbinMiddleware, "admin"), r.UsersController.Restore)
}
