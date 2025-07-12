package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"wb-L0/modules/redis"
	"wb-L0/structs"

	redis_lib "github.com/go-redis/redis/v8"
)

type RedisCache struct {
	redisConn *redis.Redis
}

func NewRedisCache(redisInstance *redis.Redis) *RedisCache {
	return &RedisCache{
		redisConn: redisInstance,
	}
}

func (c *RedisCache) PutOrder(ctx context.Context, key string, order *structs.Order) error {
	jsonData, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("marshaling error: %w", err)
	}
	return c.redisConn.Client.Set(ctx, key, jsonData, 24*time.Hour).Err()
}

func (c *RedisCache) GetOrder(ctx context.Context, key string) (*structs.Order, error) {
	val, err := c.redisConn.Client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis_lib.Nil) {
			return nil, ErrCacheMiss{
				Key: key,
			}
		}
		return nil, err
	}
	var order structs.Order
	if err := json.Unmarshal(val, &order); err != nil {
		return nil, fmt.Errorf("unmarshaling error: %w", err)
	}
	return &order, nil
}

// HealthCheck performs a health check on Redis
func (c *RedisCache) HealthCheck(ctx context.Context) error {
	return c.redisConn.Client.Ping(ctx).Err()
}
