package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/tirtahakimpambudhi/restful_api/internal/entity"
	"time"

	"github.com/phuslu/log"
	"github.com/redis/go-redis/v9"
)

// CacheRepository defines the interface for cache operations.
type CacheRepository[T any] interface {
	GetFromCache(ctx context.Context, key string) ([]T, error)      // Fetch data from cache
	SetToCache(ctx context.Context, key string, entities []T) error // Store data in cache
	DeleteToCacheByRegexKey(ctx context.Context, key string) error  // Delete cache entries by regex key
	DeleteToCache(ctx context.Context, key string) error            // Delete a specific cache entry
}

// CacheRepositoryImpl implements the CacheRepository interface using Redis.
type CacheRepositoryImpl[T any] struct {
	Cache  *redis.Client // Redis client for cache operations
	Logger *log.Logger   // Logger for logging cache operations
}

// NewCacheRepository creates a new CacheRepositoryImpl instance.
func NewCacheRepository[T any](cache *redis.Client, logger *log.Logger) *CacheRepositoryImpl[T] {
	return &CacheRepositoryImpl[T]{Cache: cache, Logger: logger}
}

func NewUserCacheRepository(cache *redis.Client, logger *log.Logger) *CacheRepositoryImpl[*entity.Users] {
	return NewCacheRepository[*entity.Users](cache, logger)
}

// GetFromCache retrieves data from the cache using the provided key.
func (r CacheRepositoryImpl[T]) GetFromCache(ctx context.Context, key string) ([]T, error) {
	r.Logger.Info().Msgf("Fetching from cache: %s", key) // Log cache fetch attempt
	cachedData, err := r.Cache.Get(ctx, key).Result()    // Get data from cache
	if errors.Is(err, redis.Nil) {
		r.Logger.Warn().Msgf("Cache miss: %s", key) // Log cache miss
		return nil, nil                             // Return nil on cache miss
	} else if err != nil {
		r.Logger.Error().Msgf("Cache error: %v", err) // Log cache error
		return nil, err                               // Return error on cache issue
	}

	var entities []T
	if err := json.Unmarshal([]byte(cachedData), &entities); err != nil {
		r.Logger.Error().Msgf("Error unmarshalling cached data: %v", err) // Log unmarshalling error
		return nil, err                                                   // Return error on unmarshalling failure
	}

	r.Logger.Info().Msgf("Cache hit: %s", key) // Log successful cache fetch
	return entities, nil
}

// SetToCache stores data in the cache with the provided key.
func (r *CacheRepositoryImpl[T]) SetToCache(ctx context.Context, key string, entities []T) error {
	data, err := json.Marshal(entities) // Convert data to JSON
	if err != nil {
		r.Logger.Error().Msgf("Error marshalling data: %v", err) // Log marshalling error
		return err                                               // Return error on marshalling failure
	}

	err = r.Cache.Set(ctx, key, data, time.Minute*30).Err() // Set data in cache with expiration
	if err != nil {
		r.Logger.Error().Msgf("Error setting cache: %v", err) // Log cache set error
		return err                                            // Return error on cache set failure
	}

	r.Logger.Info().Msgf("Data cached: %s", key) // Log successful cache set
	return nil
}

// DeleteFromCache removes a specific cache entry using the provided key.
func (r *CacheRepositoryImpl[T]) DeleteToCache(ctx context.Context, key string) error {
	err := r.Cache.Del(ctx, key).Err() // Delete data from cache
	if err != nil {
		r.Logger.Error().Msgf("Failed to delete cache for key %s: %v", key, err) // Log cache delete error
		return err                                                               // Return error on cache delete failure
	}
	return nil
}

// DeleteToCacheByRegexKey deletes cache entries that match the provided regex key.
func (r *CacheRepositoryImpl[T]) DeleteToCacheByRegexKey(ctx context.Context, key string) error {
	keys := []string{}                              // Initialize list to hold keys
	iter := r.Cache.Scan(ctx, 0, key, 0).Iterator() // Scan for keys matching regex
	for iter.Next(ctx) {
		keys = append(keys, iter.Val()) // Collect matching keys
	}
	if err := iter.Err(); err != nil {
		r.Logger.Error().Msgf("Failed to scan keys: %v", err) // Log scan error
		return err                                            // Return error on scan failure
	}
	if len(keys) > 0 {
		if err := r.Cache.Del(ctx, keys...).Err(); err != nil {
			r.Logger.Error().Msgf("Failed to delete for keys %p: %v", keys, err) // Log delete error
			return err                                                           // Return error on delete failure
		}
	}
	r.Logger.Info().Msgf("Successfully deleted cache for keys %p", keys) // Log successful cache delete
	return nil
}
