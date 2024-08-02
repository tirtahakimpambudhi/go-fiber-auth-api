package timeout

import (
	"context"
	"os"
	"strconv"
	"time"
)

type Config struct {
	CacheTimeout      time.Duration
	DatabaseTimeout   time.Duration
	DownstreamTimeout time.Duration
}

func (config Config) CreateCacheTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, config.CacheTimeout)
}

func (config Config) CreateDatabaseTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, config.DatabaseTimeout)
}

func (config Config) CreateDownstreamTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, config.DownstreamTimeout)
}

func NewConfig() (*Config, error) {
	cacheTimeout, err := strconv.Atoi(os.Getenv("CACHE_TIMEOUT"))
	if err != nil {
		return nil, err
	}
	databaseTimeout, err := strconv.Atoi(os.Getenv("DB_TIMEOUT"))
	if err != nil {
		return nil, err
	}
	downstreamTimeout, err := strconv.Atoi(os.Getenv("DOWN_STREAM_TIMEOUT"))
	if err != nil {
		return nil, err
	}
	return &Config{
		CacheTimeout:      time.Duration(cacheTimeout) * time.Minute,
		DatabaseTimeout:   time.Duration(databaseTimeout) * time.Minute,
		DownstreamTimeout: time.Duration(downstreamTimeout) * time.Minute,
	}, nil
}
