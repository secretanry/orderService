package cache

import (
	"container/list"
	"context"
	"sync"

	"wb-L0/structs"
)

type entry struct {
	key   string
	value *structs.Order
}

type MemoryCache struct {
	capacity int
	cacheMap map[string]*list.Element
	list     *list.List
	mutex    sync.Mutex
}

// NewMemoryCache creates a new LRU cache with specified capacity
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		mutex:    sync.Mutex{},
		capacity: 10,
		cacheMap: make(map[string]*list.Element),
		list:     list.New(),
	}
}

func (c *MemoryCache) PutOrder(_ context.Context, key string, order *structs.Order) error {
	c.mutex.Lock()
	if elem, exists := c.cacheMap[key]; exists {
		elem.Value.(*entry).value = order
		c.list.MoveToFront(elem)
		return nil
	}

	newEntry := &entry{key: key, value: order}
	newElem := c.list.PushFront(newEntry)
	c.cacheMap[key] = newElem

	if c.list.Len() > c.capacity {
		tail := c.list.Back()
		if tail != nil {
			keyToRemove := tail.Value.(*entry).key
			delete(c.cacheMap, keyToRemove)
			c.list.Remove(tail)
		}
	}
	c.mutex.Unlock()
	return nil
}

func (c *MemoryCache) GetOrder(_ context.Context, key string) (*structs.Order, error) {
	c.mutex.Lock()
	elem, exists := c.cacheMap[key]
	if !exists {
		return nil, ErrCacheMiss{Key: key}
	}

	c.list.MoveToFront(elem)
	c.mutex.Unlock()
	return elem.Value.(*entry).value, nil
}
