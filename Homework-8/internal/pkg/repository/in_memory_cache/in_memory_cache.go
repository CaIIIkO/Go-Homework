package inmemorycache

import (
	"sync"
	"time"

	"homework-3/internal/pkg/repository"
)

type PvzCache struct {
	data    map[int64]cachedPvz
	cleanup time.Duration
	mu      sync.RWMutex
}

type cachedPvz struct {
	repository.Pvz
	expiry time.Time
}

var Cache *PvzCache

func NewPvzCache(cleanupInterval time.Duration) {
	cache := &PvzCache{
		data:    make(map[int64]cachedPvz),
		cleanup: cleanupInterval,
	}
	go cache.cleanupExpired()
	Cache = cache
}

func (c *PvzCache) GetPvz(id int64) (repository.Pvz, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cachedPvz, found := c.data[id]
	if !found || time.Now().After(cachedPvz.expiry) {
		return repository.Pvz{}, false
	}

	return cachedPvz.Pvz, true
}

func (c *PvzCache) SetPvz(id int64, pvz repository.Pvz, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[id] = cachedPvz{
		Pvz:    pvz,
		expiry: time.Now().Add(ttl),
	}
}

func (c *PvzCache) DeletePvz(id int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, id)
}

func (c *PvzCache) cleanupExpired() {
	ticker := time.NewTicker(c.cleanup)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for id, item := range c.data {
			if time.Now().After(item.expiry) {
				delete(c.data, id)
			}
		}
		c.mu.Unlock()
	}
}
