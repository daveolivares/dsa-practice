package linkedlist

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

// assertList checks both the contents and Len, so a counter that drifts
// out of sync with the actual nodes is caught by every test.
func assertList(t *testing.T, l *LinkedList, want []int) {
	t.Helper()
	if got := l.ToSlice(); !slices.Equal(got, want) {
		t.Errorf("ToSlice() = %v, want %v", got, want)
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
		{"head", []int{1, 2, 3}, 1, true, []int{2, 3}},
		{"middle", []int{1, 2, 3}, 2, true, []int{1, 3}},
		{"last", []int{1, 2, 3}, 3, true, []int{1, 2}},
		{"duplicates at head", []int{1, 1, 2, 3, 4, 5}, 1, true, []int{1, 2, 3, 4, 5}},
		{"duplicates after head", []int{1, 2, 2, 3}, 2, true, []int{1, 2, 3}},
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
		{"last", []int{1, 2, 3}, 3, true},
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

// Sequences catch bugs that single operations miss, such as a stale head
// after the list is emptied or a counter that drifts over many operations.
func TestSequences(t *testing.T) {
	t.Run("one element, AddFirst then DeleteLast", func(t *testing.T) {
		l := &LinkedList{}
		l.AddFirst(1)
		if v, ok := l.DeleteLast(); v != 1 || !ok {
			t.Fatalf("DeleteLast() = (%d, %t), want (1, true)", v, ok)
		}
		assertList(t, l, nil)
	})

	t.Run("reuse after emptying", func(t *testing.T) {
		l := newList(1, 2)
		l.DeleteFirst()
		l.DeleteLast()
		assertList(t, l, nil)

		l.AddLast(3)
		l.AddFirst(4)
		assertList(t, l, []int{4, 3})
	})

	t.Run("delete from empty does not change Len", func(t *testing.T) {
		l := &LinkedList{}
		l.DeleteFirst()
		l.DeleteLast()
		l.DeleteValue(1)
		assertList(t, l, nil)

		l.AddFirst(1)
		assertList(t, l, []int{1})
	})

	t.Run("mixed operations", func(t *testing.T) {
		l := &LinkedList{}
		l.AddLast(2)     // [2]
		l.AddFirst(1)    // [1 2]
		l.AddLast(3)     // [1 2 3]
		l.DeleteValue(2) // [1 3]
		l.AddLast(4)     // [1 3 4]
		l.DeleteFirst()  // [3 4]
		l.DeleteLast()   // [3]
		assertList(t, l, []int{3})
	})
}
