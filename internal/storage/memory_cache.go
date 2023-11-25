package storage

import "sync"

type CacheI interface {
	Get(key string) ([]byte, bool)
	Put(key string, data []byte)
	Len() int
}

type inmemoryCache struct {
	mapa  map[string][]byte
	mutex sync.RWMutex
}

func NewInmemoryCache() *inmemoryCache {
	return &inmemoryCache{mapa: make(map[string][]byte)}
}

func (c *inmemoryCache) Get(key string) ([]byte, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	value, ok := c.mapa[key]
	return value, ok
}

func (c *inmemoryCache) Len() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.mapa)
}

func (c *inmemoryCache) Put(key string, data []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.mapa[key] = data
}
