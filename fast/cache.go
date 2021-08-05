package fast

type fastCache struct {
	shards    []*cacheShard
	shardMask uint64
	hash      fnv64a
}

func (c *fastCache) Set(key string, value interface{}) {
	c.getShard(key).set(key, value)
}

func (c *fastCache) Get(key string) interface{} {
	return c.getShard(key).get(key)
}

func (c *fastCache) Del(key string) {
	c.getShard(key).del(key)
}

func (c *fastCache) DelOldest() {
	panic("implement me")
}

func (c *fastCache) Len() int {
	length := 0
	for _, shard := range c.shards {
		length += shard.len()
	}
	return length
}

func NewFastCache(maxEntries, shardsNum int, onEvicted func(key string, value interface{})) *fastCache {
	fastCache := &fastCache{
		hash:      newDefaultHasher(),
		shards:    make([]*cacheShard, shardsNum),
		shardMask: uint64(shardsNum - 1),
	}
	for i := 0; i < shardsNum; i++ {
		fastCache.shards[i] = newCacheShard(maxEntries, onEvicted)
	}
	return fastCache
}

// maxEntries 表示分片的最大容纳记录数
// 在具体的set 或 get 等方法中 关键点是通过key 获取对应的分片，所以定义如下方法
func (c *fastCache) getShard(key string) *cacheShard {
	hashedKey := c.hash.Sum64(key)
	return c.shards[hashedKey&c.shardMask]
}
