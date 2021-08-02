package lfu

import (
	"cache"
	"container/heap"
)

// lfu 是一个LFU cache 它不是并发安全的 最少使用
type lfu struct {
	// 缓存最大容量，单位字节
	// groupcache 使用的是最大存放entry个数
	maxBytes int
	// 当一个entry 从缓存中移除时调用该回调函数 默认nil
	// groupcache 中的key 是任意的可比较类型 value 是interface{}
	onEvicted func(key string, value interface{})

	//已使用字节数
	usedBytes int

	queue *queue
	cache map[string]*entry
}

func (l *lfu) Set(key string, value interface{}) {
	if e, ok := l.cache[key]; ok {
		l.usedBytes = l.usedBytes - cache.CalcLen(e.value) + cache.CalcLen(value)
		l.queue.update(e, value, e.weight+1)
		return
	}
	en := &entry{key: key, value: value}
	heap.Push(l.queue, en)
	l.cache[key] = en
	l.usedBytes += en.Len()
	if l.maxBytes > 0 && l.usedBytes > l.maxBytes {
		l.removeElement(heap.Pop(l.queue))
	}
}

func (l *lfu) Get(key string) interface{} {
	if e, ok := l.cache[key]; ok {
		l.queue.update(e, e.value, e.weight+1)
		return e.value
	}
	return nil
}

func (l *lfu) Del(key string) {
	if e,ok:=l.cache[key];ok {
		heap.Remove(l.queue, e.index)
		l.removeElement(e)
	}
}

func (l *lfu) DelOldest() {
	if l.queue.Len() == 0 {
		return
	}
	l.removeElement(heap.Pop(l.queue))
}

func (l *lfu) removeElement(x interface{}) {
	if x == nil {
		return
	}
	en := x.(*entry)
	delete(l.cache, en.key)
	l.usedBytes -= en.Len()
	if l.onEvicted != nil {
		l.onEvicted(en.key, en.value)
	}
}

func (l *lfu) Len() int {
	return l.queue.Len()
}

type entry struct {
	key    string
	value  interface{}
	weight int
	index  int
}

func (e *entry) Len() int {
	return cache.CalcLen(e.value) + 4 + 4
}

// queue是一个entry 切片指针
// 这里的entry 和FIFO 算法的区别是多了weight index
// weight 权重优先级
// index 表示该entry 堆(heap) 中的索引
// LFU 使用最小堆实现 通过标准库container/heap 来实现最小堆 要求queue 实现heap.interface接口
type queue []*entry

//
//type Interface interface {
//	sort.Interface
//	Push(x interface{}) // add x as element Len()
//	Pop() interface{}   // remove and return element Len() - 1.s
//}

// Len queue 实现如下
func (q queue) Len() int {
	return len(q)
}

// Less < 小堆 > 大堆
func (q queue) Less(i, j int) bool {
	return q[i].weight < q[j].weight
}

func (q queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *queue) Push(x interface{}) {
	n := len(*q)
	en := x.(*entry)
	en.index = n
	*q = append(*q, en)
}

func (q *queue) Pop() interface{} {
	old := *q
	n := len(old)
	en := old[n-1]
	old[n-1] = nil // avoid memory leak
	en.index = -1  // for safety
	*q = old[0 : n-1]
	return en
}

func (q *queue) update(en *entry, value interface{}, weight int) {
	en.value = value
	en.weight = weight
	heap.Fix(q, en.index)
}

func New(maxBytes int, onEvicted func(key string, value interface{})) cache.Cache {
	q := make(queue, 0, 1024)
	return &lfu{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		queue:     &q,
		cache:     make(map[string]*entry),
	}
}

func (l *lfu) Getq() *queue {
	return l.queue
}
