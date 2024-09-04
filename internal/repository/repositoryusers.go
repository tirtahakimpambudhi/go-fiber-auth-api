package repository

import (
	"context"
	"errors"
	"github.com/phuslu/log"
	"github.com/tirtahakimpambudhi/restful_api/internal/entity"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/request"
	"gorm.io/gorm"
)

// UsersRepository defines the methods for interacting with the users repository.
type UsersRepository interface {
	Create(ctx context.Context, entity *entity.Users) error
	Update(ctx context.Context, entity *entity.Users, id any) error
	Delete(ctx context.Context, id any) error
	Restore(ctx context.Context, id any) error
	CountById(ctx context.Context, id any) (int64, error)
	Count(ctx context.Context) (int64, error)
	ExistByKeyValue(ctx context.Context, keyvalue map[string]any) (bool, error)
	GetById(ctx context.Context, entity *entity.Users, id any) error
	GetByEmail(ctx context.Context, entity *entity.Users, email string) error
	GetAll(ctx context.Context, queryParams *request.Page) ([]*entity.Users, error)
}

// UsersRepositoryImpl implements the UsersRepository interface.
type UsersRepositoryImpl struct {
	*Repository[entity.Users]             // Embedded generic repository
	DB                        *gorm.DB    // Database connection
	Logger                    *log.Logger // Logger for logging messages
}

// NewUsersRepositoryImpl creates a new instance of UsersRepositoryImpl.
func NewUsersRepositoryImpl(DB *gorm.DB, logger *log.Logger) (*UsersRepositoryImpl, error) {
	// Check if DB or logger is nil
	if DB == nil || logger == nil {
		return nil, errors.New("DB or Logger is nil")
	}
	// Return new instance of UsersRepositoryImpl
	return &UsersRepositoryImpl{Repository: NewRepository[entity.Users](logger, DB), DB: DB, Logger: logger}, nil
}

// GetAll retrieves all users based on the query parameters.
func (repo UsersRepositoryImpl) GetAll(ctx context.Context, queryParams *request.Page) ([]*entity.Users, error) {
	// Log the start of the GetAll method
	repo.Logger.Info().Msg("Starting GetAll method")
	var users []*entity.Users

	// Prepare the query with context and default limit
	queryDB := repo.DB.Unscoped().WithContext(ctx).Model(&entity.Users{}).Limit(queryParams.Size)

	// Apply cursor pagination if Before parameter is provided
	if queryParams.Before != "" {
		repo.Logger.Info().Msgf("Applying cursor pagination with Before ID: %s", queryParams.Before)
		queryDB = queryDB.Where("id < ?", queryParams.Before)
	}
	// Apply cursor pagination if After parameter is provided
	if queryParams.After != "" {
		repo.Logger.Info().Msgf("Applying cursor pagination with After ID: %s", queryParams.After)
		queryDB = queryDB.Where("id > ?", queryParams.After)
	}
	// Apply default ordering by ID in ascending order
	repo.Logger.Info().Msg("Applying default order by ID ascending")
	queryDB = queryDB.Order("id ASC")

	// Execute the query and fetch users
	err := queryDB.Find(&users).Error
	if err != nil {
		// Log error if fetching users fails
		repo.Logger.Error().Msgf("Error occurred while fetching users: %v", err)
		return nil, err
	}

	// Log successful retrieval of users
	repo.Logger.Info().Msgf("Successfully retrieved %d users", len(users))
	return users, nil
}

// ExistByKeyValue checks if a user exists based on key-value conditions.
func (repo UsersRepositoryImpl) ExistByKeyValue(ctx context.Context, keyvalue map[string]any) (bool, error) {
	var countResult int64
	query := repo.DB.WithContext(ctx).Model(&entity.Users{})
	// Apply where conditions for each key-value pair
	for key, value := range keyvalue {
		repo.Logger.Info().Msgf("Applying where conditions column %s value %v", key, value)
		query = query.Where(key+" = ?", value)
	}
	// Count the number of records matching the conditions
	err := query.Count(&countResult).Error
	repo.Logger.Info().Msgf("Count result: %d", countResult)
	return countResult == 1, err
}

func (repo UsersRepositoryImpl) GetByEmail(ctx context.Context, users *entity.Users, email string) error {
	repo.Logger.Info().Msgf("Retrieving entity by Email: %v", email)
	// Retrieve entity by ID from the database
	err := repo.DB.Model(&entity.Users{}).WithContext(ctx).Where("email = ?", email).Take(users).Error
	if err != nil {
		// Log error if retrieval failed.
		repo.Logger.Error().Msgf("Failed to retrieve entity by Email: %v", err)
		return err
	}
	// Log success if entity retrieval was successful.
	repo.Logger.Info().Msg("Entity successfully retrieved")
	return nil
}

func (repo UsersRepositoryImpl) Restore(ctx context.Context, id any) error {
	repo.Logger.Info().Msgf("Retrieving entity by ID: %v", id)
	tx, closeTx := repo.Transaction(ctx) // Start transaction
	defer closeTx(tx)                    // Ensure transaction is committed or rolled back
	err := tx.Unscoped().Model(&entity.Users{}).WithContext(ctx).Where("id = ?", id).Not("deleted_at", 0).Update("deleted_at", 0).Error
	if err != nil {
		//Rollback Transaction
		tx.Rollback()
		// Log error if retrieval failed.
		repo.Logger.Error().Msgf("Failed to restore entity by ID: %v", err)
		return err
	}
	// Log success if entity retrieval was successful.
	repo.Logger.Info().Msg("Entity successfully restore")
	return nil
}
