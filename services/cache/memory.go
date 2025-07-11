package cache

import (
	"context"

	"wb-L0/structs"
)

type MemoryCache struct {
	Data map[string]*structs.Order
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		Data: make(map[string]*structs.Order),
	}
}

func (c *MemoryCache) PutOrder(_ context.Context, key string, order *structs.Order) error {
	c.Data[key] = order
	return nil
}

func (c *MemoryCache) GetOrder(_ context.Context, key string) (*structs.Order, error) {
	value, ok := c.Data[key]
	if !ok {
		return nil, ErrCacheMiss{
			Key: key,
		}
	}
	return value, nil
}
