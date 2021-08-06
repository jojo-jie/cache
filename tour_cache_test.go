package cache_test

import (
	"cache"
	"cache/fast"
	"cache/lru"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	is2 "github.com/matryer/is"
)

func TestTourCacheGet(t *testing.T) {
	db := map[string]string{
		"key1": "val1",
		"key2": "val2",
		"key3": "val3",
		"key4": "val4",
		"key5": "val5",
	}

	getter := cache.GetFunc(func(key string) interface{} {
		log.Println("[From DB ] find key ", key)
		if val, ok := db[key]; ok {
			return val
		}
		return nil
	})

	tourCache := cache.NewTourCache(getter, lru.New(0, nil))
	is := is2.New(t)
	var wg sync.WaitGroup
	for k, v := range db {
		wg.Add(1)
		go func(k, v string) {
			defer wg.Done()
			is.Equal(tourCache.Get(k), v)
			is.Equal(tourCache.Get(k), v)
		}(k, v)
	}
	wg.Wait()

	is.Equal(tourCache.Get("unknown"), nil)
	is.Equal(tourCache.Get("unknown"), nil)

	is.Equal(tourCache.Stat().NGet, 12)
	is.Equal(tourCache.Stat().NHit, 5)

}

//  -run=none 只运行benchmark -count 输出次数
//  go test -bench=. -benchmem -count=2 -run=none. | tee old.txt
//1. BenchmarkFast-8 表测试的函数名，-8 表示GOMAXPROCS（线程数）的值为8
//
//2. 10000000 表一共执行了一千万次，即B.N的值
//
//3. 107 ns/op表平均每次操作花费了107纳秒
//
//4. 16 B/op 表每次操作申请了16Byte的内存申请
//
//5. 2 allocs/op 表每次操作申请了2次内存
//1.  参数-bench，它指明要测试的函数；点字符意思是测试当前所有以Benchmark为前缀函数
//
//2.  参数-benchmem，性能测试的时候显示测试函数的内存分配大小，内存分配次数的统计信息
//
//3. 参数-count n,运行测试和性能多少此，默认一次
//b.ResetTimer()，代表重置计时为0，以调用时的时刻作为重新计时的开始
func BenchmarkFast(b *testing.B) {
	cache := fast.NewFastCache(b.N, 100, nil)
	rand.Seed(time.Now().Unix())
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(100)
		counter := 0
		for pb.Next() {
			cache.Set("s:"+strconv.Itoa(id), id)
			counter++
		}
	})
}
