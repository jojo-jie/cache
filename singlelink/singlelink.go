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
func (l *LinkedList) Invert(n uint, s *ListNode) {
	if l.head == nil || l.head.next == nil {
		fmt.Println("出错了!!!")
		return
	}
	//1 2 3
	var pre, temp *ListNode
	cur := s
	var i uint
	for cur != nil {
		i++
		temp = cur.next
		cur.next = pre
		pre = cur
		cur = temp
		if n > 0 && i == n {
			s.next = temp
			break
		}
	}
	l.head.next = pre
}

func (l *LinkedList) InvertMN(m, n uint) {
	if l.head == nil || l.head.next == nil {
		fmt.Println("出错了!!!")
		return
	}
	cur := l.head.next
	var i uint
	var pre, temp, s, e, o *ListNode
	for cur != nil {
		i++
		temp = cur.next
		if i+1 == m {
			s = cur
		}
		if i >= m && i <= n {
			if i == m {
				e = cur
			}
			cur.next = pre
			pre = cur
			cur = temp
			if i == n {
				o = temp
			}
		}

		cur = temp
		if i > n {
			s.next = pre
			e.next = o
		}
	}
}

func Reverse(head *ListNode) *ListNode {
	// 递归找到链表最后一个节点
	if head == nil || head.next == nil {
		return head
	} else {
		newhead := Reverse(head.next)
		head.next.next = head
		head.next = nil
		return newhead
	}
}

var successor *ListNode = nil

func ReverseN(head *ListNode, n uint) *ListNode {
	if n == 1 {
		// 记录第n+1 个节点
		successor = head.next
		return head
	}
	last := ReverseN(head.next, n-1)
	head.next.next = head
	head.next = successor
	return last
}

func ReverseBetween(head *ListNode, m, n uint) *ListNode {
	// base case
	if m == 1 {
		return ReverseN(head, n)
	}
	// 前进到反转起点触发
	head.next = ReverseBetween(head, m-1, n-1)
	return head
}



func MergeTwoLink(l1, l2 *ListNode) *ListNode {
	if l1 == nil {
		return l2
	}
	if l1 == nil {
		return l1
	}

	var n *ListNode
	if l1.data.(int) < l2.data.(int) {
		n = l1
		n.next = MergeTwoLink(l1.next, l2)
	} else {
		n = l2
		n.next = MergeTwoLink(l1, l2.next)
	}
	return n
}
