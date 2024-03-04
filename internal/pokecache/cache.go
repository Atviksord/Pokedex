package main

import (
	"sync"
	"time"
)

type Cache struct {
	cacheMap map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache { // returns a new instance (new memory address everytime) of the Cache struct using pointers

	cache := &Cache{
		cacheMap: make(map[string]cacheEntry),
		interval: interval,
	}
	cache.reapLoop()
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	c.cacheMap[key] = cacheEntry{val: val, createdAt: time.Now()}

	c.mu.Unlock()
}
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if cacheEntry, exists := c.cacheMap[key]; exists {

		return cacheEntry.val, true

	} else {

		return nil, false
	}

}
func (c *Cache) reapLoop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	ticker := time.NewTicker(c.interval)
	for range ticker.C {

	}

}

func main() {
	cache := NewCache(5 * time.Minute)

}
