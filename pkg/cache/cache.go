package cache

import "sync"

type Cache[K comparable, V any] struct {
	entries map[K]V
	mu      sync.RWMutex
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.get(key)
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.set(key, value)
}

func (c *Cache[K, V]) GetFunc(key K, fn func() (V, error)) (V, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if v, ok := c.get(key); ok {
		return v, nil
	}
	v, err := fn()
	if err != nil {
		return v, err
	}
	c.set(key, v)
	return v, nil
}

func (c *Cache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[K]V)
}

func (c *Cache[K, V]) get(key K) (V, bool) {
	v, ok := c.entries[key]
	return v, ok
}

func (c *Cache[K, V]) set(key K, value V) {
	if c.entries == nil {
		c.entries = make(map[K]V)
	}
	c.entries[key] = value
}
