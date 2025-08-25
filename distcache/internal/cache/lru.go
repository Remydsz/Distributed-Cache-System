package cache

import (
	"container/list"
	"sync"
	"time"
)

type entry struct {
	key string
	val []byte
	exp time.Time      // zero => no TTL
	el  *list.Element  // node in LRU list
}

type LRU struct {
	cap int
	ll  *list.List                 // front = MRU, back = LRU
	m   map[string]*entry
	mu  sync.RWMutex
}

func NewLRU(capacity int) *LRU {
	return &LRU{
		cap: capacity,
		ll:  list.New(),
		m:   make(map[string]*entry),
	}
}

func (c *LRU) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.m[key]
	if !ok {
		return nil, false
	}
	// TTL check
	if !e.exp.IsZero() && time.Now().After(e.exp) {
		c.remove(e)
		return nil, false
	}
	c.ll.MoveToFront(e.el)
	return e.val, true
}

func (c *LRU) Set(key string, val []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if e, ok := c.m[key]; ok {
		e.val = val
		e.exp = expireAt(ttl)
		c.ll.MoveToFront(e.el)
		return
	}
	e := &entry{key: key, val: val, exp: expireAt(ttl)}
	e.el = c.ll.PushFront(e)
	c.m[key] = e

	if len(c.m) > c.cap {
		c.evictOne()
	}
}

func (c *LRU) Delete(key string) {
	c.mu.Lock()
	if e, ok := c.m[key]; ok {
		c.remove(e)
	}
	c.mu.Unlock()
}

func (c *LRU) evictOne() {
	if back := c.ll.Back(); back != nil {
		c.remove(back.Value.(*entry))
	}
}

func (c *LRU) remove(e *entry) {
	c.ll.Remove(e.el)
	delete(c.m, e.key)
}

func expireAt(ttl time.Duration) time.Time {
	if ttl <= 0 {
		return time.Time{}
	}
	return time.Now().Add(ttl)
}
