package fast

import (
	"container/list"
	"sync"
)

type cacheShard struct {
	locker sync.RWMutex

	// 最大存放 entry 个数
	maxEntries int
	onEvicted  func(key string, value interface{})
	ll         *list.List
	cache      map[string]*list.Element
}

func (c *cacheShard) set(key string, value interface{}) {
	c.locker.Lock()
	defer c.locker.Unlock()
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToBack(e)
		en := e.Value.(*entry)
		en.value = value
		return
	}
	en := &entry{key: key, value: value}
	e := c.ll.PushBack(en)
	c.cache[key] = e
	if c.maxEntries > 0 && c.ll.Len() > c.maxEntries {
		c.removeElement(c.ll.Front())
	}
}

func (c *cacheShard) get(key string) interface{} {
	c.locker.RLock()
	defer c.locker.RUnlock()
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToBack(e)
		return e.Value.(*entry).value
	}
	return nil
}

func (c *cacheShard) del(key string) {
	c.locker.Lock()
	defer c.locker.Unlock()
	if e, ok := c.cache[key]; ok {
		c.removeElement(e)
	}
}

func (c *cacheShard) delOldest() {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.removeElement(c.ll.Front())
}

func (c *cacheShard) len() int {
	c.locker.RLock()
	defer c.locker.RUnlock()
	return c.ll.Len()
}

func (c *cacheShard) removeElement(e *list.Element) {
	if e == nil {
		return
	}
	c.ll.Remove(e)
	en := e.Value.(*entry)
	delete(c.cache, en.key)
	if c.onEvicted != nil {
		c.onEvicted(en.key, en.value)
	}
}

// 创建一个新的 cacheShard ，如果maxBytes是0 ，则表示没有容量限制
func newCacheShard(maxEntries int, onEvicted func(key string, value interface{})) *cacheShard {
	return &cacheShard{
		maxEntries: maxEntries,
		onEvicted:  onEvicted,
		ll:         list.New(),
		cache:      make(map[string]*list.Element),
	}
}

type entry struct {
	key   string
	value interface{}
}
