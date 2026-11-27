// Package stack implements a last-in, first-out (LIFO) stack of ints,
// backed by a slice whose end is the top of the stack.
package stack

// Stack is a LIFO stack. The zero value is an empty stack ready to use.
type Stack struct {
	items []int
}

// Push adds value to the top of the stack.
//
// Time: amortized O(1). When the slice is full, append copies every element
// into a larger array, but the capacity grows geometrically, so the average
// cost per push stays constant.
func (s *Stack) Push(value int) {
	s.items = append(s.items, value)
}

// Pop removes the top value and returns it with true. If the stack is
// empty, it returns 0 and false.
//
// Time: O(1).
func (s *Stack) Pop() (int, bool) {
	if len(s.items) == 0 {
		return 0, false
	}
	top := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return top, true
}

// Peek returns the top value with true, without removing it. If the stack
// is empty, it returns 0 and false.
//
// Time: O(1).
func (s *Stack) Peek() (int, bool) {
	if len(s.items) == 0 {
		return 0, false
	}
	return s.items[len(s.items)-1], true
}

// IsEmpty reports whether the stack has no values.
//
// Time: O(1).
func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

// Len returns the number of values on the stack.
//
// Time: O(1).
func (s *Stack) Len() int {
	return len(s.items)
}
