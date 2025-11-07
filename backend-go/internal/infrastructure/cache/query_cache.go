package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
)

// QueryResultCache implementa caché para resultados de consultas
type QueryResultCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewQueryResultCache(client *redis.Client, ttl time.Duration) *QueryResultCache {
	return &QueryResultCache{
		client: client,
		ttl:    ttl,
	}
}

// Set almacena un resultado en caché
func (c *QueryResultCache) Set(ctx context.Context, key string, result *entities.QueryResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal query result: %w", err)
	}

	return c.client.Set(ctx, c.makeKey(key), data, c.ttl).Err()
}

// Get obtiene un resultado de la caché
func (c *QueryResultCache) Get(ctx context.Context, key string) (*entities.QueryResult, error) {
	data, err := c.client.Get(ctx, c.makeKey(key)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Caché miss
		}
		return nil, err
	}

	var result entities.QueryResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal query result: %w", err)
	}

	return &result, nil
}

// Delete elimina un resultado de la caché
func (c *QueryResultCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, c.makeKey(key)).Err()
}

// makeKey genera una clave de caché para una consulta
func (c *QueryResultCache) makeKey(key string) string {
	return fmt.Sprintf("query_result:%s", key)
}

// GenerateKey genera una clave única para una consulta
func (c *QueryResultCache) GenerateKey(userID uint, sql string, database string) string {
	return fmt.Sprintf("%d:%s:%s", userID, database, sql)
}

// Clear limpia todas las claves que coinciden con un patrón
func (c *QueryResultCache) Clear(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, fmt.Sprintf("query_result:%s*", pattern), 0).Iterator()
	for iter.Next(ctx) {
		if err := c.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}
