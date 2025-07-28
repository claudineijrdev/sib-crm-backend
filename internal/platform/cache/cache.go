package cache

import (
	"context"
	"time"
)

type CacheService interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}

type cacheService struct {
	// Implementação específica (Redis, Memory, etc.)
	client interface{}
}

func NewCacheService(client interface{}) CacheService {
	return &cacheService{
		client: client,
	}
}

func (c *cacheService) Get(ctx context.Context, key string) (interface{}, error) {
	// Implementação real aqui
	return nil, nil
}

func (c *cacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// Implementação real aqui
	return nil
}

func (c *cacheService) Delete(ctx context.Context, key string) error {
	// Implementação real aqui
	return nil
}

func (c *cacheService) Exists(ctx context.Context, key string) (bool, error) {
	// Implementação real aqui
	return false, nil
} 