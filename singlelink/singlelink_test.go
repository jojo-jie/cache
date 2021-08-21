package singlelink

import "testing"

var link *LinkedList
var n = 3

func init()  {
	link=NewLinkedList()
	cur:=link.head
	for i := 1; i < n; i++ {
		link.InsertAfter(cur, i)
		cur = cur.next
	}
}

// 链表反转
func TestLinked(t *testing.T) {
	t.Log(link.Print())
	link.Invert()
	t.Log(link.Print())
}
