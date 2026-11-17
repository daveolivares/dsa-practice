// Package linkedlist implements a singly linked list of ints.
//
// The list keeps a pointer to its first node and a count of its nodes, but no
// tail pointer, so operations at the end of the list take O(n) time, where n
// is the number of elements.
package linkedlist

// node is one element of the list. next is nil for the last node.
type node struct {
	value int
	next  *node
}

// LinkedList is a singly linked list. The zero value is an empty list ready
// to use.
type LinkedList struct {
	head   *node
	length int
}

// AddFirst inserts value at the front of the list.
//
// Time: O(1).
func (l *LinkedList) AddFirst(value int) {
	newNode := &node{value: value, next: l.head}
	l.head = newNode
	l.length++
}

// AddLast appends value to the end of the list.
//
// Time: O(n), because the list has no tail pointer and must be walked to
// find the last node.
func (l *LinkedList) AddLast(value int) {
	newNode := &node{value: value, next: nil}
	if l.head == nil {
		l.head = newNode
	} else {
		curr := l.head
		for curr.next != nil {
			curr = curr.next
		}
		curr.next = newNode
	}
	l.length++
}

// Contains reports whether value is in the list.
//
// Time: O(n).
func (l *LinkedList) Contains(value int) bool {
	curr := l.head
	for curr != nil {
		if curr.value == value {
			return true
		}
		curr = curr.next
	}
	return false
}

// DeleteFirst removes the first element and returns it with true. If the
// list is empty, it returns 0 and false.
//
// Time: O(1).
func (l *LinkedList) DeleteFirst() (int, bool) {
	if l.head == nil {
		return 0, false
	}
	v := l.head.value
	l.head = l.head.next
	l.length--
	return v, true
}

// DeleteLast removes the last element and returns it with true. If the list
// is empty, it returns 0 and false.
//
// Time: O(n), because the list must be walked to find the node before the
// last one.
func (l *LinkedList) DeleteLast() (int, bool) {
	if l.head == nil {
		return 0, false
	}

	if l.head.next == nil {
		v := l.head.value
		l.head = nil
		l.length--
		return v, true
	}

	curr := l.head
	for curr.next.next != nil {
		curr = curr.next
	}
	v := curr.next.value
	curr.next = nil
	l.length--
	return v, true
}

// DeleteValue removes the first element equal to value and reports whether
// one was removed. Later duplicates are left in place.
//
// Time: O(n).
func (l *LinkedList) DeleteValue(value int) bool {
	if l.head == nil {
		return false
	}

	if l.head.value == value {
		l.DeleteFirst()
		return true
	}

	prev, curr := l.head, l.head.next
	for curr != nil {
		if curr.value == value {
			prev.next = curr.next
			l.length--
			return true
		}
		prev = curr
		curr = curr.next
	}
	return false
}

// Len returns the number of elements in the list.
//
// Time: O(1).
func (l *LinkedList) Len() int {
	return l.length
}

// ToSlice returns the elements in order, from first to last. The slice is a
// copy, so changing it does not change the list.
//
// Time: O(n).
func (l *LinkedList) ToSlice() []int {
	slc := make([]int, 0, l.length)
	curr := l.head
	for curr != nil {
		slc = append(slc, curr.value)
		curr = curr.next
	}
	return slc
}
