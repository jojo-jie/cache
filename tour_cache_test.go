package cache

import (
	"cache/lru"
	is2 "github.com/matryer/is"
	"log"
	"sync"
	"testing"
)

func TestTourCacheGet(t *testing.T) {
	db:=map[string]string{
		"key1":"val1",
		"key2":"val2",
		"key3":"val3",
		"key4":"val4",
		"key5":"val5",
	}

	getter:= GetFunc(func(key string) interface{} {
		log.Println("[From DB ] find key ",key)
		if val,ok:=db[key];ok {
			return val
		}
		return nil
	})

	tourCache:= NewTourCache(getter,lru.New(0,nil))
	is:=is2.New(t)
	var wg sync.WaitGroup
	for k, v := range db {
		wg.Add(1)
		go func(k ,v string) {
			is.Equal(tourCache.Get(k), v)
			is.Equal(tourCache.Get(k), v)
		}(k,v)
	}
	wg.Wait()

	is.Equal(tourCache.Get("unknown"),nil)
	is.Equal(tourCache.Get("unknown"),nil)

	is.Equal(tourCache.Stat().NGet, 10)
	is.Equal(tourCache.Stat().NHit, 4)

}
