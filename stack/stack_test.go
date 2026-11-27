package stack

import "testing"

// assertLen checks that Len and IsEmpty agree with the expected size.
func assertLen(t *testing.T, s *Stack, want int) {
	t.Helper()
	if got := s.Len(); got != want {
		t.Errorf("Len() = %d, want %d", got, want)
	}
	if got := s.IsEmpty(); got != (want == 0) {
		t.Errorf("IsEmpty() = %t, want %t", got, want == 0)
	}
}

// assertPop pops once and checks the returned value and ok flag.
func assertPop(t *testing.T, s *Stack, wantV int, wantOK bool) {
	t.Helper()
	if v, ok := s.Pop(); v != wantV || ok != wantOK {
		t.Errorf("Pop() = (%d, %t), want (%d, %t)", v, ok, wantV, wantOK)
	}
}

// assertPeek peeks once and checks the returned value and ok flag.
func assertPeek(t *testing.T, s *Stack, wantV int, wantOK bool) {
	t.Helper()
	if v, ok := s.Peek(); v != wantV || ok != wantOK {
		t.Errorf("Peek() = (%d, %t), want (%d, %t)", v, ok, wantV, wantOK)
	}
}

func TestEmptyStack(t *testing.T) {
	s := &Stack{}
	assertLen(t, s, 0)
	assertPop(t, s, 0, false)
	assertPeek(t, s, 0, false)
	assertLen(t, s, 0) // Pop and Peek on an empty stack change nothing
}

func TestPush(t *testing.T) {
	t.Run("onto empty stack", func(t *testing.T) {
		s := &Stack{}
		s.Push(1)
		assertLen(t, s, 1)
		assertPeek(t, s, 1, true)
	})

	t.Run("new value becomes the top", func(t *testing.T) {
		s := &Stack{}
		s.Push(1)
		s.Push(2)
		assertLen(t, s, 2)
		assertPeek(t, s, 2, true)
	})
}

func TestPop(t *testing.T) {
	t.Run("returns the top and removes it", func(t *testing.T) {
		s := &Stack{}
		s.Push(1)
		s.Push(2)
		assertPop(t, s, 2, true)
		assertLen(t, s, 1)
		assertPeek(t, s, 1, true)
	})

	t.Run("last element empties the stack", func(t *testing.T) {
		s := &Stack{}
		s.Push(1)
		assertPop(t, s, 1, true)
		assertLen(t, s, 0)
		assertPop(t, s, 0, false)
	})

	t.Run("zero is a real value", func(t *testing.T) {
		s := &Stack{}
		s.Push(0)
		assertPeek(t, s, 0, true)
		assertPop(t, s, 0, true)
	})
}

// Added by mentor.
func TestPeekDoesNotRemove(t *testing.T) {
	s := &Stack{}
	s.Push(1)
	s.Push(2)
	assertPeek(t, s, 2, true)
	assertPeek(t, s, 2, true)
	assertLen(t, s, 2)
	assertPop(t, s, 2, true)
}

func TestLIFOOrder(t *testing.T) {
	t.Run("push all then pop all", func(t *testing.T) {
		s := &Stack{}
		for i := 1; i <= 5; i++ {
			s.Push(i)
		}
		assertLen(t, s, 5)
		for want := 5; want >= 1; want-- {
			assertPop(t, s, want, true)
		}
		assertLen(t, s, 0)
		assertPop(t, s, 0, false)
	})

	// Added by mentor.
	t.Run("interleaved push and pop", func(t *testing.T) {
		s := &Stack{}
		s.Push(1)
		s.Push(2)
		assertPop(t, s, 2, true)
		s.Push(3)
		assertPop(t, s, 3, true)
		assertPop(t, s, 1, true)
		assertLen(t, s, 0)
	})
}

// Added by mentor.
func TestReuseAfterEmptying(t *testing.T) {
	s := &Stack{}
	s.Push(1)
	s.Pop()
	assertLen(t, s, 0)

	s.Push(2)
	assertLen(t, s, 1)
	assertPeek(t, s, 2, true)
	assertPop(t, s, 2, true)
}

// Added by mentor. Enough pushes to make append grow the slice many times.
func TestManyPushes(t *testing.T) {
	const n = 10000
	s := &Stack{}
	for i := 0; i < n; i++ {
		s.Push(i)
	}
	assertLen(t, s, n)
	for want := n - 1; want >= 0; want-- {
		if v, ok := s.Pop(); v != want || !ok {
			t.Fatalf("Pop() = (%d, %t), want (%d, true)", v, ok, want)
		}
	}
	assertLen(t, s, 0)
}
