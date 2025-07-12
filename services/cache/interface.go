package cache

import (
	"context"

	"wb-L0/structs"
)

var cacheInstance Cache

type Cache interface {
	GetOrder(context.Context, string) (*structs.Order, error)
	PutOrder(context.Context, string, *structs.Order) error
	HealthCheck(context.Context) error
}

func SetCache(cache Cache) {
	cacheInstance = cache
}

func GetCache() Cache {
	return cacheInstance
}
