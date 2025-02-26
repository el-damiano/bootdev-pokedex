package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu      sync.Mutex
	entries map[string]cacheEntry
}

type cacheEntry struct {
	creationTime time.Time
	val          []byte
}

func NewCache(interval time.Duration) *Cache {
	newCache := Cache{
		entries: map[string]cacheEntry{},
	}
	go newCache.reapLoop(interval)
	return &newCache
}

func (c *Cache) Add(key string, value []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	newEntry := cacheEntry{
		creationTime: time.Now(),
		val:          value,
	}

	c.entries[key] = newEntry
	return nil
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entries[key]
	if !ok {
		return nil, ok
	}
	return entry.val, ok
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reap(interval)
	}
}

func (c *Cache) reap(last time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, value := range c.entries {
		if time.Since(value.creationTime) > last {
			delete(c.entries, key)
		}
	}
}
