package pokecache

import (
	"sync"
	"time"
	"fmt"
)

type Cache struct {
	Entries	map[string]CacheEntry
	ReapInterval time.Duration
	reaperQuitChannel chan struct{}
	Mu	sync.Mutex
}

type CacheEntry struct {
	CreatedAt	time.Time
	Val	[]byte
}

func NewCache(interval time.Duration) *Cache {
	quitChan := make(chan struct{})
	returnCache := &Cache {
		Entries: make(map[string]CacheEntry),
		ReapInterval: interval,
		reaperQuitChannel:	quitChan,
	}
	StartReaper(interval, quitChan, returnCache)
	return returnCache
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

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
}

func (c *Cache) PrintStatus() {
	fmt.Println("Current state of the cache:")
	for key := range c.Entries {
		fmt.Println(key)
	}
}

func (c *Cache) ReapEntries() {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	for key, entry := range c.Entries {
		if time.Since(entry.CreatedAt) > c.ReapInterval {
			delete(c.Entries, key)
		}
	}
}

func (c *Cache) StopReaper() {
	close(c.reaperQuitChannel)
}

func StartReaper(interval time.Duration, quit chan struct{}, cache *Cache) {
	fmt.Println("Starting the cache reaper")
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				cache.ReapEntries()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}