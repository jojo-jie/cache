package lru

import (
	is2 "github.com/matryer/is"
	"testing"
	"unsafe"
)

func TestOnEvicted(t *testing.T) {
	is := is2.New(t)
	keys := make([]string, 0, 8)
	onEvicted := func(key string, value interface{}) {
		keys = append(keys, key)
	}
	cache := New(16, onEvicted)
	cache.Set("k1", 1)
	cache.Set("k2", 2)
	cache.Get("k1")
	cache.Set("k3", 3)
	cache.Get("k1")
	cache.Set("k4", 4)
	expected := []string{"k2", "k3"}
	is.Equal(expected, keys)
}

//unsafe.Sizeof 返回数据类型大小 空结构体大小==0
//unsafe.Alignof 返回类型对齐系数
//内存对其
//对于结构体的各个成员，第一个成员位于偏移为0的位置，结构体第一个成员的偏移量(offset)为0，以后每个成员相对于结构体首地址的offset都是该成员大小与有效对齐值中较小那个的整数倍，如有需要编译器会在成员之间加上填充字节。
//除了结构成员需要对齐，结构本身也需要对齐，结构的长度必须是编译器默认的对齐长度和成员中最长类型中最小的数据大小的倍数对齐。

type Demo struct {
	A int32
	B []int32
	C string
	D bool
	E struct{}
}

type test1 struct {
	// 大小 对齐值
	a bool   // 1 1
	b int32  // 4 4
	c string // 16 8
}

type test2 struct {
	a int32  // 4 4
	b string // 16 8
	c bool   // 1 1
}

type test3 struct {
	b string // 16 8
	c bool   // 1 1
	a int32  // 4 4
}

type test4 struct {
	b []int32
	a struct{} //struct{} 在最后会被填充对齐到前一个字段的对齐系数
}

func TestStruct(t *testing.T) {
	t4 := test4{}
	t.Log(unsafe.Alignof(t4.b))
	t.Log(unsafe.Sizeof(t4))
}

func TestSize(t *testing.T) {
	var s Demo
	t.Log(unsafe.Alignof(s.A))
	t.Log(unsafe.Sizeof(s.A))
	t.Log(unsafe.Sizeof(s.B))
	t.Log(unsafe.Sizeof(s.C))
	t.Log(unsafe.Alignof(s.E))
	// 0-3 4-7 8-31 32 47 48 56

	// 0 1-3 4-7 8->23 => 24
	var t1 test1
	t.Log(unsafe.Sizeof(t1))

	// 0-3 4-7 8->23 24 31 => 32
	var t2 test2
	t.Log(unsafe.Sizeof(t2))

	// 0-15 16 17-19 20-23=>24
	var t3 test3
	t.Log(unsafe.Sizeof(t3))

}
