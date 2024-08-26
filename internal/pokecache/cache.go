package pokecache

import (
	"sync"
	"time"
)

// CacheEntry represents a single entry in the cache
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// Cache is the main structure that holds the cache map and a mutex for synchronization
type Cache struct {
	data     map[string]cacheEntry
	mutex    sync.Mutex
	interval time.Duration
	reapDone chan bool
}

// NewCache creates a new Cache with the given interval for reaping old entries
func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		data:     make(map[string]cacheEntry),
		interval: interval,
		reapDone: make(chan bool),
	}
	go cache.reapLoop()
	return cache
}

// Add adds a new entry to the cache
func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

// Get retrieves an entry from the cache. It returns the value and a boolean indicating if the key was found.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, found := c.data[key]
	if !found {
		return nil, false
	}

	return entry.val, true
}

// reapLoop is a method that runs in a goroutine and removes old entries from the cache at intervals
func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.reap()
		case <-c.reapDone:
			return
		}
	}
}

// reap removes any entries from the cache that are older than the cache's interval
func (c *Cache) reap() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	for key, entry := range c.data {
		if now.Sub(entry.createdAt) > c.interval {
			delete(c.data, key)
		}
	}
}

// StopReap stops the reapLoop goroutine when the cache is no longer needed
func (c *Cache) StopReap() {
	c.reapDone <- true
}
