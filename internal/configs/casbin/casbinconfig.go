package casbinconfig

import (
	"errors"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gofiber/fiber/v2"

	"github.com/tirtahakimpambudhi/restful_api/internal/configs"
	pathhelper "github.com/tirtahakimpambudhi/restful_api/pkg/helper/path"
	"gorm.io/gorm"
)

// NewCasbin initializes a new Casbin middleware for Fiber with a GORM database.
func NewCasbin(db *gorm.DB, lookup func(ctx *fiber.Ctx) string) (*casbin.Enforcer, error) {
	// Check if the database connection is nil
	if db == nil {
		return nil, errors.New("db is nil") // Return error if the database is nil
	}

	var casbinInstance Casbin
	// Load Casbin configuration from the environment
	err := configs.GetConfig().Load(&casbinInstance)
	if err != nil {
		return nil, err // Return error if loading configuration fails
	}

	// Create a new GORM adapter for Casbin
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err // Return error if creating adapter fails
	}

	// Create a new Casbin enforcer with the model path and adapter
	enforce, err := casbin.NewEnforcer(pathhelper.AddWorkdirToSomePath(casbinInstance.ModelPath, casbinInstance.ModelName), adapter)
	if err != nil {
		return nil, err // Return error if creating enforcer fails
	}

	// Return a new Casbin middleware instance for Fiber
	return enforce, nil
}

// Casbin holds the configuration for Casbin.
type Casbin struct {
	ModelPath string `env:"MODEL_PATH" envDefault:"resource/model"` // Path to the Casbin model file
	ModelName string `env:"MODEL_FILENAME" envDefault:"rbac_model.conf"`
}
