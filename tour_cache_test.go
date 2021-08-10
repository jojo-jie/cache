package cache_test

import (
	"cache"
	"cache/fast"
	"cache/lru"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"runtime"
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

var datas []string

// 采样方式
// runtime/pprof 采集程序 指定区块的运行数据进行分析
// net/http/pprof 基于http server运行，并且可以采集运行时的数据进行分析
// go test 通过运行测试用例 指定所需要标识进行采集
// 使用模式
// report generation 报告生成
// 交互式终端使用
// web界面 http://127.0.0.1:port/debug/pprof  +debug=1 可直接在浏览器中访问 无denbug参数 直接下载profile 文件
// full goroutine stack map
// allocs 查看过去所有内存分配的样本，访问路径为$HOST/debug/pprof/allocs
// block 查看导致阻塞同步的堆栈跟踪，访问路径为$HOST/debug/pprof/block
// cmdline 当前程序命令行的完整调用路径
// goroutine 查看当前所有运行的goroutines 堆栈跟踪 访问路径为$HOST/debug/pprof/goroutine
// heap 查看活动对象的内存分配情况
// mutex 查看导致互斥锁的竞争持有者的堆栈跟踪 访问路径为$HOST/debug/pprof/mutex
// profile 默认进行30s cpu Profiling 会得到一个分析用的profile 文件 访问路径为$HOST/debug/pprof/profile
// threadcreate 查看创建os线程的堆栈跟踪 访问路径为$HOST/debug/pprof/threadcreate

func init() {
	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)
}

func TestHttpPProf(t *testing.T) {
	go func() {
		for {
			t.Logf("len: %d", Add("tour-book"))
			time.Sleep(time.Millisecond * 10)
		}
	}()
	go func() {
		_ = http.ListenAndServe(":4444", nil)
	}()
	var m sync.Mutex
	var datas = make(map[int]struct{})
	for i := 0; i < 999; i++ {
		go func(i int) {
			m.Lock()
			defer m.Unlock()
			datas[i] = struct{}{}
		}(i)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/tt", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello pprof"))
	})
	s := &http.Server{
		Addr:    ":8896",
		Handler: mux,
	}
	_ = s.ListenAndServe()
}

func Add(str string) int {
	data := []byte(str)
	datas = append(datas, string(data))
	return len(data)
}
