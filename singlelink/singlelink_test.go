package singlelink

import (
	"fmt"
	"testing"
)

var link *LinkedList
var n = 10

func init() {
	link = NewLinkedList()
	cur := link.head
	for i := 1; i < n; i++ {
		link.InsertAfter(cur, i)
		cur = cur.next
	}
}

// 链表反转
func TestLinked(t *testing.T) {
	t.Log(link.Print())
	/*link.Invert()
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
