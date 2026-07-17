package counter

import (
	"fmt"
	"sync"
)

type Counter struct {
	mu     sync.RWMutex
	values map[string]int
}

func New() *Counter {
	return &Counter{values: make(map[string]int)}
}

func (c *Counter) Add(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[key]++
}

func (c *Counter) Get(key string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.values[key]
}

func (c *Counter) Summary(key string) string {
	return fmt.Sprintf("key=%s count=%d", key, c.Get(key))
}
