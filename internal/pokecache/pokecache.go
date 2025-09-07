package pokecache

import (
	"sync"
	"time"
	"fmt"
)

type Cache struct {
	Entries	map[string]CacheEntry
	Mu	sync.Mutex
}

type CacheEntry struct {
	CreatedAt	time.Time
	Val	[]byte
}

func NewCache(interval time.Duration) *Cache {
	returnCache := &Cache {
		Entries: make(map[string]CacheEntry),
	}
	return returnCache
}

func (c *Cache) Get(key string) ([]byte, bool) {
	entries, ok := c.Entries[key]
	if !ok {
		return nil, false
	}
	return entries.Val, true
}

func (c *Cache) Add(key string, data []byte) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	var entry CacheEntry
	entry.CreatedAt = time.Now()
	entry.Val = data
	c.Entries[key] = entry

	c.PrintStatus()
}

func (c *Cache) PrintStatus() {
	fmt.Println("Current state of the cache:")
	for key := range c.Entries {
		fmt.Println(key)
	}
}