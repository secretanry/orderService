package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"wb-L0/modules/config"
	"wb-L0/modules/graceful"
)

type Redis struct {
	Client *redis.Client
}

func (r *Redis) Init(_ chan error) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.GetConfig().RedisHost, config.GetConfig().RedisPort),
		Username:     "default",
		Password:     config.GetConfig().RedisPass,
		DB:           config.GetConfig().RedisDatabase,
		PoolSize:     20,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
	})

	if err := pingRedis(rdb); err != nil {
		return fmt.Errorf("redis connection failed: %v", err)
	}
	r.Client = rdb
	return nil
}

func (r *Redis) SuccessfulMessage() string {
	return "Redis successfully initialized"
}

func (r *Redis) Shutdown(_ context.Context) error {
	return r.Client.Close()
}

func pingRedis(rdb *redis.Client) error {
	ctx, cancel := context.WithTimeout(graceful.GetContext(), 3*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	return err
}
