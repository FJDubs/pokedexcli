package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	Entries map[string]cacheEntry
	mu      sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	newCache := &Cache{
		Entries: make(map[string]cacheEntry),
	}
	go newCache.reapLoop(interval)
	return newCache
}

func (c *Cache) Add(s string, data []byte) {
	c.mu.Lock()
	c.Entries[s] = cacheEntry{time.Now(), data}
	c.mu.Unlock()
}

func (c *Cache) Get(s string) ([]byte, bool) {
	c.mu.Lock()
	if entry, ok := c.Entries[s]; !ok {
		c.mu.Unlock()
		return nil, false
	} else {
		c.mu.Unlock()
		return entry.val, true
	}
}

func (c *Cache) reapLoop(interval time.Duration) {
	tkr := time.NewTicker(interval)
	for {
		t := <-tkr.C
		c.mu.Lock()
		for name, entry := range c.Entries {
			if interval < t.Sub(entry.createdAt) {
				delete(c.Entries, name)

			}
		}
		c.mu.Unlock()

	}

}
