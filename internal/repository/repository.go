package repository

import (
	"context"

	"github.com/phuslu/log"
	"gorm.io/gorm"
)

// Repository is a generic repository struct for handling database operations.
type Repository[T any] struct {
	Logger *log.Logger // Logger for logging operations.
	DB     *gorm.DB    // Database connection.
}

// NewRepository creates a new instance of Repository with the provided logger and database connection.
func NewRepository[T any](logger *log.Logger, db *gorm.DB) *Repository[T] {
	return &Repository[T]{Logger: logger, DB: db}
}

// CommitOrRollback handles transaction commit or rollback based on whether an error occurred.
func (r *Repository[T]) CommitOrRollback(tx *gorm.DB) {
	if err := recover(); err != nil {
		// Log error and rollback transaction if an error occurred.
		r.Logger.Error().Msgf("Failed to commit Database : %v", err)
		tx.Rollback()
	} else {
		// Commit transaction if no error occurred.
		tx.Commit()
	}
}

// Transaction starts a new database transaction and returns the transaction and a function to commit or rollback.
func (r *Repository[T]) Transaction(ctx context.Context) (*gorm.DB, func(tx *gorm.DB)) {
	tx := r.DB.Begin().WithContext(ctx).Model(new(T))
	return tx, r.CommitOrRollback
}

// Create attempts to create a new entity in the database.
func (r *Repository[T]) Create(ctx context.Context, entity *T) error {
	r.Logger.Info().Msg("Attempting to create a new entity")
	tx, closeTx := r.Transaction(ctx) // Start transaction
	defer closeTx(tx)                 // Ensure transaction is committed or rolled back

	// Create entity in the database
	err := tx.Create(entity).Error
	if err != nil {
		// Log error and rollback transaction if creation failed.
		tx.Rollback()
		r.Logger.Error().Msgf("Failed to create entity: %v", err)
		return err
	}

	// Log success if entity creation was successful.
	r.Logger.Info().Msg("Entity successfully created")
	return nil
}

// Update attempts to update an existing entity in the database.
func (r *Repository[T]) Update(ctx context.Context, entity *T, id any) error {
	r.Logger.Info().Msg("Attempting to update an entity")
	tx, closeTx := r.Transaction(ctx) // Start transaction
	defer closeTx(tx)                 // Ensure transaction is committed or rolled back

	// Update entity in the database
	err := tx.Where("id = ?", id).Updates(entity).Error
	if err != nil {
		// Log error and rollback transaction if update failed.
		tx.Rollback()
		r.Logger.Error().Msgf("Failed to update entity: %v", err)
		return err
	}

	// Log success if entity update was successful.
	r.Logger.Info().Msg("Entity successfully updated")
	return nil
}

// Delete attempts to delete an entity from the database.
func (r *Repository[T]) Delete(ctx context.Context, id any) error {
	r.Logger.Info().Msg("Attempting to delete an entity")
	tx, closeTx := r.Transaction(ctx) // Start transaction
	defer closeTx(tx)                 // Ensure transaction is committed or rolled back

	// Delete entity from the database
	err := tx.Where("id = ?", id).Delete(new(T)).Error
	if err != nil {
		// Log error and rollback transaction if deletion failed.
		tx.Rollback()
		r.Logger.Error().Msgf("Failed to delete entity: %v", err)
		return err
	}

	// Log success if entity deletion was successful.
	r.Logger.Info().Msg("Entity successfully deleted")
	return nil
}

// CountById counts the number of entities with the specified ID.
func (r Repository[T]) CountById(ctx context.Context, id any) (int64, error) {
	r.Logger.Info().Msgf("Counting entities with ID: %v", id)
	var total int64
	// Count entities by ID in the database
	err := r.DB.Model(new(T)).WithContext(ctx).Where("id = ?", id).Count(&total).Error
	if err != nil {
		// Log error if counting failed.
		r.Logger.Error().Msgf("Failed to count entities by ID: %v", err)
		return 0, err
	}
	// Log success with the total count of entities.
	r.Logger.Info().Msgf("Total entities counted: %d", total)
	return total, nil
}

// Count counts the total number of entities in the database.
func (r Repository[T]) Count(ctx context.Context) (int64, error) {
	r.Logger.Info().Msg("Counting entities")
	var total int64
	// Count total entities in the database
	err := r.DB.Model(new(T)).WithContext(ctx).Count(&total).Error
	if err != nil {
		// Log error if counting failed.
		r.Logger.Error().Msg("Failed to count entities")
		return 0, err
	}
	// Log success with the total count of entities.
	r.Logger.Info().Msgf("Total entities counted: %d", total)
	return total, nil
}

// GetById retrieves an entity by its ID from the database.
func (r Repository[T]) GetById(ctx context.Context, entity *T, id any) error {
	r.Logger.Info().Msgf("Retrieving entity by ID: %v", id)
	// Retrieve entity by ID from the database
	err := r.DB.Model(new(T)).WithContext(ctx).Where("id = ?", id).Take(entity).Error
	if err != nil {
		// Log error if retrieval failed.
		r.Logger.Error().Msgf("Failed to retrieve entity by ID: %v", err)
		return err
	}
	// Log success if entity retrieval was successful.
	r.Logger.Info().Msg("Entity successfully retrieved")
	return nil
}
