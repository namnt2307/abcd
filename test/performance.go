package performance

import (
	"sync"
)

type RCache struct {
	sync.RWMutex
	cache map[string]string
}

type LCache struct {
	sync.RWMutex
	cache map[string]string
}

func (c RCache) Get(key string) string {
	c.RLock()
	v := c.cache[key]
	c.RUnlock()
	return v
}

func (c RCache) Set(key string) {
	c.Lock()
	c.cache = map[string]string{}
	c.cache[key] = key
	c.Unlock()
}

func (c LCache) Get(key string) string {
	c.Lock()
	v := c.cache[key]
	c.Unlock()
	return v
}

func (c LCache) Set(key string) {
	c.Lock()
	c.cache = map[string]string{}
	c.cache[key] = key
	c.Unlock()
}