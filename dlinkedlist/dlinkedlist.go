// Package dlinkedlist implements a doubly linked list of ints.
//
// Each node links to both its next and previous node, and the list keeps
// pointers to both ends, so adding or removing at either end takes O(1)
// time. n below is the number of elements.
package dlinkedlist

// node is one element of the list. prev is nil for the first node and next
// is nil for the last node.
type node struct {
	value int
	prev  *node
	next  *node
}

// LinkedList is a doubly linked list. The zero value is an empty list ready
// to use.
type LinkedList struct {
	head   *node
	tail   *node
	length int
}

// AddFirst inserts value at the front of the list.
//
// Time: O(1).
func (l *LinkedList) AddFirst(value int) {
	newNode := &node{value: value, next: l.head}
	if l.head == nil {
		l.tail = newNode
	} else {
		l.head.prev = newNode
	}
	l.head = newNode
	l.length++
}

// AddLast appends value to the end of the list.
//
// Time: O(1), using the tail pointer.
func (l *LinkedList) AddLast(value int) {
	newNode := &node{value: value, prev: l.tail}
	if l.tail == nil {
		l.head = newNode
	} else {
		l.tail.next = newNode
	}
	l.tail = newNode
	l.length++
}

// Contains reports whether value is in the list.
//
// Time: O(n).
func (l *LinkedList) Contains(value int) bool {
	curr := l.tail
	for curr != nil {
		if curr.value == value {
			return true
		}
		curr = curr.prev
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
	if l.head.next == nil {
		l.head = nil
		l.tail = nil
	} else {
		l.head = l.head.next
		l.head.prev = nil
	}
	l.length--
	return v, true
}

// DeleteLast removes the last element and returns it with true. If the list
// is empty, it returns 0 and false.
//
// Time: O(1), because the last node's prev link leads straight to the new
// last node.
func (l *LinkedList) DeleteLast() (int, bool) {
	if l.tail == nil {
		return 0, false
	}

	v := l.tail.value
	if l.tail.prev == nil {
		l.head = nil
		l.tail = nil
	} else {
		l.tail = l.tail.prev
		l.tail.next = nil
	}
	l.length--
	return v, true
}

// DeleteValue removes the first element, counting from the front, that is
// equal to value, and reports whether one was removed. Later duplicates are
// left in place.
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

	curr := l.head.next
	for curr != nil && curr.next != nil {
		if curr.value == value {
			curr.prev.next = curr.next
			curr.next.prev = curr.prev
			l.length--
			return true
		}
		curr = curr.next
	}

	if l.tail.value == value {
		l.DeleteLast()
		return true
	}
	return false
}

// Len returns the number of elements in the list.
//
// Time: O(1).
func (l *LinkedList) Len() int {
	return l.length
}

// Reverse reverses the order of the list in place by swapping every node's
// prev and next links, then swapping head and tail. No nodes are created.
//
// Time: O(n). Space: O(1).
func (l *LinkedList) Reverse() {
	curr := l.head
	for curr != nil {
		tempNext := curr.next
		curr.next = curr.prev
		curr.prev = tempNext
		curr = curr.prev
	}
	tempHead := l.head
	l.head = l.tail
	l.tail = tempHead
}

// ToSlice returns the elements from first to last. The slice is a copy, so
// changing it does not change the list.
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

// ToSliceReverse returns the elements from last to first, by following prev
// links from the tail. The slice is a copy.
//
// Time: O(n).
func (l *LinkedList) ToSliceReverse() []int {
	slc := make([]int, 0, l.length)
	curr := l.tail
	for curr != nil {
		slc = append(slc, curr.value)
		curr = curr.prev
	}
	return slc
}
