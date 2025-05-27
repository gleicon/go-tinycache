package datautils

import (
	"container/list"
	"sync"
)

type entry struct{ key, val string }

type LRUCache struct {
	cap  int
	list *list.List
	data map[string]*list.Element
	mu   sync.Mutex
}

func NewLRU(cap int) *LRUCache {
	return &LRUCache{
		cap:  cap,
		list: list.New(),
		data: make(map[string]*list.Element),
	}
}

func (c *LRUCache) Get(k string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.data[k]; ok {
		c.list.MoveToFront(e)
		return e.Value.(entry).val, true
	}
	return "", false
}

func (c *LRUCache) Set(k, v string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.data[k]; ok {
		c.list.MoveToFront(e)
		e.Value = entry{k, v}
		return
	}
	if c.list.Len() >= c.cap {
		old := c.list.Back()
		c.list.Remove(old)
		delete(c.data, old.Value.(entry).key)
	}
	e := c.list.PushFront(entry{k, v})
	c.data[k] = e
}
func (c *LRUCache) Delete(k string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.data[k]; ok {
		c.list.Remove(e)
		delete(c.data, k)
	}
}
func (c *LRUCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.list.Init()
	c.data = make(map[string]*list.Element)
}
func (c *LRUCache) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.list.Len()
}
