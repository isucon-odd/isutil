package cache

import (
	"sync"
	"time"
)

type cacheItem[T any] struct {
	value     T
	expiredAt *time.Time
}

type Cache[T any] struct {
	cache map[string]cacheItem[T]
	mutex *sync.RWMutex
}

func NewCache[T any]() *Cache[T] {
	return &Cache[T]{
		cache: make(map[string]cacheItem[T]),
		mutex: &sync.RWMutex{},
	}
}

func (c *Cache[T]) Set(key string, value T) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache[key] = cacheItem[T]{
		value:     value,
		expiredAt: nil,
	}
}

func (c *Cache[T]) SetWithExpiration(key string, value T, expiration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	expired := time.Now().Add(expiration)
	expiredAt := &expired

	c.cache[key] = cacheItem[T]{
		value:     value,
		expiredAt: expiredAt,
	}
}

func (c *Cache[T]) Get(key string) (T, bool) {
	var defaultValue T

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, ok := c.cache[key]
	if !ok {
		return defaultValue, false
	}
	return item.value, true
}

func (c *Cache[T]) GetAndDeleteExpired(key string) (T, bool) {
	var defaultValue T

	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, ok := c.cache[key]
	if !ok {
		return defaultValue, false
	}

	if item.expiredAt != nil && time.Now().After(*item.expiredAt) {
		delete(c.cache, key)
		return defaultValue, false
	}

	return item.value, true
}

func (c *Cache[T]) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.cache, key)
}
