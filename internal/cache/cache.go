package cache

import (
	"sync"

	"github.com/sparxfort1ano/order-service/internal/repository"
)

type OrderCache interface {
	// Check if cache contains the order
	Get(id string) (*repository.Order, bool)
	// Put order into cache
	Set(order *repository.Order)
}

type CacheMap map[string]*repository.Order

type Cache struct {
	orderMap CacheMap
	mu       sync.RWMutex
}

func (c *Cache) Get(id string) (result *repository.Order, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result, ok = c.orderMap[id]
	return
}

func (c *Cache) Set(order *repository.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.orderMap) >= 1000 {
		for k := range c.orderMap {
			delete(c.orderMap, k)
			break
		}
	}
	c.orderMap[order.OrderUid] = order
}

// Init Cache
func NewOrderCache() *Cache {
	return &Cache{
		orderMap: make(CacheMap, 1000),
	}
}
