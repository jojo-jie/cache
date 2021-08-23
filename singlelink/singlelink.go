package singlelink

import "fmt"

type ListNode struct {
	data interface{}
	next *ListNode
}

type LinkedList struct {
	head   *ListNode
	length uint
}

func NewListNode(v interface{}) *ListNode {
	return &ListNode{v, nil}
}

func (this *ListNode) GetNext() *ListNode {
	return this.next
}

func (this *ListNode) GetValue() interface{} {
	return this.data
}

func NewLinkedList() *LinkedList {
	return &LinkedList{NewListNode(0), 0}
}

//							n  o
// InsertAfter 某个节点后面 m->m->p->m->e
func (l *LinkedList) InsertAfter(p *ListNode, v interface{}) bool {
	if p == nil {
		return false
	}
	newNode := NewListNode(v)
	newNode.next = p.next
	p.next = newNode
	l.length++
	return true
}

func (l *LinkedList) InsertBefore(p *ListNode, v interface{}) bool {
	if p == nil || p == l.head {
		return false
	}

	cur := l.head.next
	pre := l.head
	for nil != cur {
		if cur == p {
			break
		}
		pre = cur
		cur = cur.next
	}
	newNode := NewListNode(v)
	newNode.next = cur
	pre.next = newNode
	l.length++
	return true
}

// InsertToHead 在链表头部插入节点
func (l *LinkedList) InsertToHead(v interface{}) bool {
	return l.InsertAfter(l.head, v)
}

func (l *LinkedList) InsertToTail(v interface{}) bool {
	cur := l.head
	for nil != cur {
		cur = cur.next
	}
	return l.InsertAfter(cur, v)
}

// FindByIndex 通过索引查找节点
func (l *LinkedList) FindByIndex(index uint) *ListNode {
	if index > l.length {
		return nil
	}

	cur := l.head.next
	var i uint
	for ; i < index; i++ {
		cur = cur.next
	}
	return cur
}

func (l *LinkedList) DeleteNode(p *ListNode) bool {
	if nil == p {
		return false
	}
	cur := l.head.next
	pre := l.head
	for nil != cur {
		if cur == p {
			break
		}
		pre = cur
		cur = cur.next
	}
	pre.next = p.next
	p = nil
	l.length--
	return true
}

// Print 打印链表
func (l *LinkedList) Print() string {
	cur := l.head.next
	var format string
	for nil != cur {
		format += fmt.Sprintf("%+v", cur.GetValue())
		cur = cur.next
		if nil != cur {
			format += "->"
		}
	}
	return format
}

// Invert 非递归链表反转
func (l *LinkedList) Invert() {
	if l.head == nil || l.head.next == nil {
		fmt.Println("出错了!!!")
		return
	}
	//1 2 3
	var pre,temp *ListNode
	cur := l.head.next
	for cur != nil {
		temp = cur.next
		cur.next = pre
		pre = cur
		cur = temp
	}
	l.head.next = pre
}

func ReverseList(head *ListNode) *ListNode {
	// 递归找到链表最后一个节点
	if head == nil || head.next == nil {
		return head
	} else {
		newhead:=ReverseList(head.next)
		head.next.next = head
		head.next = nil
		return newhead
	}
}

