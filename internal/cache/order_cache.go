package cache

import (
	"sync"

	"github.com/sparxfort1ano/order-service/internal/repository"
)

// Кэш заказов в памяти
type OrderCache struct {
	mu sync.RWMutex                // блокировка для потокобезопасности
	m  map[string]repository.Order // словарь: order_uid -> заказ
}

// Создание нового пустого кэша
func New() *OrderCache { return &OrderCache{m: make(map[string]repository.Order)} }

// Добавляем или обновляем заказ в кэше
func (c *OrderCache) Set(o repository.Order) {
	c.mu.Lock()
	c.m[o.OrderUID] = o
	c.mu.Unlock()
}

// Получаем заказ из кэша по id
func (c *OrderCache) Get(id string) (repository.Order, bool) {
	c.mu.RLock()
	o, ok := c.m[id]
	c.mu.RUnlock()
	return o, ok
}

// Массовое заполнение кэша списком заказов
func (c *OrderCache) Warm(list []repository.Order) {
	c.mu.Lock()
	for _, o := range list {
		c.m[o.OrderUID] = o
	}
	c.mu.Unlock()
}
