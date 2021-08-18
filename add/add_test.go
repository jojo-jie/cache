package add

import (
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

//结构体多字段原子操作

type Person struct {
	name string
	age  int
}

var p Person
var m sync.Mutex

func update(name string, age int) {
	m.Lock()
	defer m.Unlock()
	p.name = name
	time.Sleep(time.Millisecond * 200)
	p.age = age
}

var a atomic.Value

func atomicUpdate(name string, age int) {
	p1 := &Person{}
	p1.name = name
	time.Sleep(time.Millisecond * 200)
	p1.age = age
	a.Store(p1)
}

func TestBatch(t *testing.T) {
	// A WaitGroup must not be copied after first use.
	// sync 结构体 传指针
	wg := sync.WaitGroup{}
	n := 10
	wg.Add(n)
	for i := 0; i < n; i++ {
		name, age := strconv.Itoa(i), i
		go func(name string, age int) {
			defer wg.Done()
			atomicUpdate(name, age)
		}(name, age)
	}
	wg.Wait()

	//Store() 参数是局部变量(一块全新的内存) Load() 没有数据拷贝
	//store value 不是指针会多一层struct 数据拷贝
	t.Logf("name:%s age:%d", a.Load().(*Person).name, a.Load().(*Person).age)
}
