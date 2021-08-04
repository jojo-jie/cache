package cache

type Getter interface {
	Get(key string) interface{}
}

type GetFunc func(key string) interface{}

func (f GetFunc) Get(key string) interface{} {
	return f(key)
}

type TourCache struct {
	mainCache *safeCache
	getter    Getter
}

func (t *TourCache) Get(key string) interface{} {
	val := t.mainCache.Get(key)
	if val != nil {
		return val
	}
	if t.getter != nil {
		val = t.getter.Get(key)
		if val == nil {
			return nil
		}
		t.mainCache.Set(key, val)
		return val
	}
	return nil
}

func (t *TourCache) Stat() *Stat {
	return t.mainCache.Stat()
}

func NewTourCache(getter Getter, cache Cache) *TourCache {
	return &TourCache{
		mainCache: newSafeCache(cache),
		getter:    getter,
	}
}
