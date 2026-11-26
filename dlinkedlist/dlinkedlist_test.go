package dlinkedlist

import (
	"slices"
	"testing"
)

// newList builds a list by calling AddLast for each value in order.
func newList(vals ...int) *LinkedList {
	l := &LinkedList{}
	for _, v := range vals {
		l.AddLast(v)
	}
	return l
}

// assertList checks the forward walk, the backward walk, and Len. The
// backward walk follows prev links from tail, so a broken prev link, a stale
// tail, or a head whose prev isn't nil makes it disagree with the forward walk.
func assertList(t *testing.T, l *LinkedList, want []int) {
	t.Helper()
	if got := l.ToSlice(); !slices.Equal(got, want) {
		t.Errorf("ToSlice() = %v, want %v", got, want)
	}
	wantRev := slices.Clone(want)
	slices.Reverse(wantRev)
	if got := l.ToSliceReverse(); !slices.Equal(got, wantRev) {
		t.Errorf("ToSliceReverse() = %v, want %v", got, wantRev)
	}
	if got := l.Len(); got != len(want) {
		t.Errorf("Len() = %d, want %d", got, len(want))
	}
}

func TestAddFirst(t *testing.T) {
	tests := []struct {
		name  string
		start []int
		value int
		want  []int
	}{
		{"empty list", nil, 1, []int{1}},
		{"one element", []int{2}, 1, []int{1, 2}},
		{"several elements", []int{2, 3}, 1, []int{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newList(tt.start...)
			l.AddFirst(tt.value)
			assertList(t, l, tt.want)
		})
	}
}

func TestAddLast(t *testing.T) {
	tests := []struct {
		name  string
		start []int
		value int
		want  []int
	}{
		{"empty list", nil, 1, []int{1}},
		{"one element", []int{1}, 2, []int{1, 2}},
		{"several elements", []int{1, 2}, 3, []int{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newList(tt.start...)
			l.AddLast(tt.value)
			assertList(t, l, tt.want)
		})
	}
}

func TestDeleteFirst(t *testing.T) {
	tests := []struct {
		name   string
		start  []int
		wantV  int
		wantOK bool
		want   []int
	}{
		{"empty list", nil, 0, false, nil},
		{"one element", []int{1}, 1, true, nil},
		{"several elements", []int{1, 2, 3}, 1, true, []int{2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newList(tt.start...)
			v, ok := l.DeleteFirst()
			if v != tt.wantV || ok != tt.wantOK {
				t.Errorf("DeleteFirst() = (%d, %t), want (%d, %t)", v, ok, tt.wantV, tt.wantOK)
			}
			assertList(t, l, tt.want)
		})
	}
}

func TestDeleteLast(t *testing.T) {
	tests := []struct {
		name   string
		start  []int
		wantV  int
		wantOK bool
		want   []int
	}{
		{"empty list", nil, 0, false, nil},
		{"one element", []int{1}, 1, true, nil},
		{"two elements", []int{1, 2}, 2, true, []int{1}},
		{"several elements", []int{1, 2, 3}, 3, true, []int{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newList(tt.start...)
			v, ok := l.DeleteLast()
			if v != tt.wantV || ok != tt.wantOK {
				t.Errorf("DeleteLast() = (%d, %t), want (%d, %t)", v, ok, tt.wantV, tt.wantOK)
			}
			assertList(t, l, tt.want)
		})
	}
}

func TestDeleteValue(t *testing.T) {
	tests := []struct {
		name  string
		start []int
		value int
		ok    bool
		want  []int
	}{
		{"empty list", nil, 1, false, nil},
		{"value not in list", []int{1, 2, 3}, 4, false, []int{1, 2, 3}},
		{"only element", []int{1}, 1, true, nil},
		{"one element, value not in list", []int{1}, 2, false, []int{1}},
		{"head", []int{1, 2, 3}, 1, true, []int{2, 3}},
		{"middle", []int{1, 2, 3}, 2, true, []int{1, 3}},
		{"tail", []int{1, 2, 3}, 3, true, []int{1, 2}},
		{"duplicates in middle", []int{1, 2, 2, 3}, 2, true, []int{1, 2, 3}},
		{"duplicates at head", []int{1, 1, 2, 3}, 1, true, []int{1, 2, 3}},
		{"duplicates at tail", []int{1, 2, 3, 3}, 3, true, []int{1, 2, 3}},
		{"non-adjacent duplicates", []int{1, 2, 3, 2}, 2, true, []int{1, 3, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newList(tt.start...)
			if got := l.DeleteValue(tt.value); got != tt.ok {
				t.Errorf("DeleteValue(%d) = %t, want %t", tt.value, got, tt.ok)
			}
			assertList(t, l, tt.want)
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name  string
		start []int
		value int
		want  bool
	}{
		{"empty list", nil, 1, false},
		{"value not in list", []int{1, 2, 3}, 4, false},
		{"head", []int{1, 2, 3}, 1, true},
		{"middle", []int{1, 2, 3}, 2, true},
		{"tail", []int{1, 2, 3}, 3, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newList(tt.start...)
			if got := l.Contains(tt.value); got != tt.want {
				t.Errorf("Contains(%d) = %t, want %t", tt.value, got, tt.want)
			}
			assertList(t, l, tt.start) // Contains must not change the list
		})
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name  string
		start []int
		want  []int
	}{
		{"empty list", nil, nil},
		{"one element", []int{1}, []int{1}},
		{"two elements", []int{1, 2}, []int{2, 1}},
		{"several elements", []int{1, 2, 3, 4}, []int{4, 3, 2, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newList(tt.start...)
			l.Reverse()
			assertList(t, l, tt.want)
		})
	}

	t.Run("AddLast after reverse", func(t *testing.T) {
		l := newList(1, 2, 3)
		l.Reverse()
		l.AddLast(4)
		assertList(t, l, []int{3, 2, 1, 4})
	})
}

// Sequences catch state left behind by earlier calls, such as a tail that
// still points at a removed node after the list is emptied.
func TestSequences(t *testing.T) {
	t.Run("empty with DeleteFirst, then AddLast", func(t *testing.T) {
		l := newList(1)
		l.DeleteFirst()
		assertList(t, l, nil)

		l.AddLast(2)
		assertList(t, l, []int{2})
	})

	t.Run("empty with DeleteLast, then AddFirst", func(t *testing.T) {
		l := newList(1)
		l.DeleteLast()
		assertList(t, l, nil)

		l.AddFirst(2)
		assertList(t, l, []int{2})
	})

	t.Run("empty with DeleteValue, then add at both ends", func(t *testing.T) {
		l := newList(1)
		l.DeleteValue(1)
		assertList(t, l, nil)

		l.AddLast(2)
		l.AddFirst(1)
		assertList(t, l, []int{1, 2})
	})

	t.Run("delete from empty does not change Len", func(t *testing.T) {
		l := &LinkedList{}
		l.DeleteFirst()
		l.DeleteLast()
		l.DeleteValue(1)
		assertList(t, l, nil)
	})

	t.Run("mixed operations", func(t *testing.T) {
		l := &LinkedList{}
		l.AddLast(2)     // [2]
		l.AddFirst(1)    // [1 2]
		l.AddLast(3)     // [1 2 3]
		l.DeleteValue(2) // [1 3]
		l.Reverse()      // [3 1]
		l.AddFirst(4)    // [4 3 1]
		l.DeleteLast()   // [4 3]
		assertList(t, l, []int{4, 3})
	})
}
