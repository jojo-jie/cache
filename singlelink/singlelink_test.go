package singlelink

import (
	"fmt"
	"testing"
)

var link *LinkedList
var n = 10
var m = make(map[int]int, 10)

func init() {
	link = NewLinkedList()
	cur := link.head
	for i := 1; i < n; i++ {
		m[i] = i
	}
	for k, v := range m {
		if k == 0 {
			continue
		}
		link.InsertAfter(cur, v)
		cur = cur.next
	}

}

// 链表反转
func TestLinked(t *testing.T) {
	t.Log(link.Print())
	/*link.Invert(6, link.head.next)
	link.InvertMN(3, 6)
	t.Log(link.Print())*/
	n := ReverseN(link.head.next, 4)
	s := ""
	for {
		if n == nil {
			break
		}
		s += fmt.Sprintf("%+v", n.GetValue())
		n = n.next
		if n != nil {
			s += "->"
		}
	}
	t.Log(s)
}
