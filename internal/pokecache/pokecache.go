package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu       sync.Mutex
	entries  map[string]cacheEntry
	interval int
}

type cacheEntry struct {
	creationTime time.Time
	val          []byte
}

func NewCache(interval time.Duration) *Cache {
	newCache := Cache{
		interval: int(interval),
		entries:  map[string]cacheEntry{},
	}
	go newCache.reapLoop()
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

func (c *Cache) reapLoop() {
	for {
		c.mu.Lock()
		for key, value := range c.entries {
			if time.Since(value.creationTime) > time.Duration(c.interval) {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
		time.Sleep(time.Duration(c.interval))
	}
}
