# Stack

Implementation: [`stack.go`](stack.go) · Tests: [`stack_test.go`](stack_test.go) · Related: [Singly Linked List](../linkedlist/LINKED_LIST.md)

## Core idea

A **stack** is a collection where the last value added is the first one removed. This rule is called **LIFO**: last in, first out. A stack of plates works the same way: you add a plate to the top, and you take one from the top.

A stack has a small set of operations:

| Operation | Meaning |
|---|---|
| `Push(v)` | Add `v` to the top. |
| `Pop()` | Remove the top value and return it. |
| `Peek()` | Return the top value without removing it. |
| `IsEmpty()`, `Len()` | Report whether the stack is empty, and how many values it holds. |

Values below the top can't be reached directly. That restriction is what makes a stack simple and fast.

This implementation stores the values in a Go slice, with the **end of the slice as the top**:

```text
Push(1), Push(2), Push(3):

index:  0   1   2
items: [1,  2,  3]
                ▲
               top = items[len(items)-1]

Pop() returns 3:

items: [1,  2]
            ▲
           top
```

## Invariant

> **The top of the stack is `items[len(items)-1]`, and `len(items)` is the number of values on the stack.**

Below the top, values are in the order they were pushed. Because the slice tracks its own length, there is no separate counter that could drift out of sync with the contents.

## Complexity

| Operation | Time | Why |
|---|---|---|
| `Push` | **amortized O(1)** | `append` writes one value to the end. Occasionally it must also grow the slice, as explained below. |
| `Pop` | O(1) | Reads the last value and shortens the slice. Nothing is copied. |
| `Peek` | O(1) | Reads the last value. |
| `IsEmpty`, `Len` | O(1) | Read the slice's length. |

**Space:** O(n) for n values, plus some unused capacity left over from growth.

### Why `Push` is amortized O(1)

A slice is a view onto an underlying array of fixed size, called its **capacity**. While there is spare capacity, `append` just writes the new value into the next slot, which is O(1).

When the array is full, `append` allocates a larger array and copies every existing value into it. That single push costs O(n). Go grows small slices by doubling their capacity, and large ones by a smaller factor. Running 2,000 pushes on an empty slice produced these capacities:

```text
4  8  16  32  64  128  256  512  848  1280  1792  2560
```

Because the capacity grows by a multiple each time, rather than by a fixed amount, the expensive copies become rarer as the slice gets bigger. Over those 2,000 pushes, about 4,900 values were copied in total, roughly 2.5 copies per push on average. The total copying work stays proportional to the number of pushes, so the **average** cost of a push is constant.

That is what **amortized** O(1) means: a single operation can occasionally be expensive, but the average over many operations is O(1).

If the capacity grew by a fixed amount instead, such as 10 more slots each time, the number of copies would grow with the square of the number of pushes, and `Push` would be O(n) on average.

## Implementation

```go
type Stack struct {
    items []int
}

func (s *Stack) Push(value int) {
    s.items = append(s.items, value)
}

func (s *Stack) Pop() (int, bool) {
    if len(s.items) == 0 {
        return 0, false
    }
    top := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return top, true
}

func (s *Stack) Peek() (int, bool) {
    if len(s.items) == 0 {
        return 0, false
    }
    return s.items[len(s.items)-1], true
}
```

`IsEmpty` and `Len` both return values based on `len(s.items)`.

### Why the top is the end of the slice

Adding or removing at the **end** of a slice moves nothing else. Adding or removing at the **front**, index 0, means every other value has to shift over by one position, which costs O(n) per operation. A stack only ever touches one end, so it should be the cheap end.

### Removing the top value

`Pop` shortens the slice with a **slice expression**: `s.items[:len(s.items)-1]` is a view of every element except the last. It doesn't copy anything, which is why `Pop` is O(1). `Pop` must read the top value before shortening the slice, because afterwards the value is outside the slice.

### Reporting an empty stack

`Pop` and `Peek` return `(int, bool)`. On an empty stack they return `(0, false)`.

A single "empty" return value such as `-1` or `0` can't work, because a stack of `int`s can hold any `int`, including that one. The boolean says whether a value was returned, so `(0, true)` (the top value is 0) and `(0, false)` (the stack is empty) stay distinct. This is the same "comma ok" pattern Go uses for map lookups: `v, ok := m[key]`.

### No constructor needed

Go struct fields can't have default values, so `items` starts as a `nil` slice. A `nil` slice has length 0 and works with `append`, so `&Stack{}` is already a valid empty stack.

When the expected size is known in advance, a constructor can reserve capacity up front and avoid the copies made during growth:

```go
func New(capacity int) *Stack {
    return &Stack{items: make([]int, 0, capacity)}
}
```

`make([]int, 0, capacity)` creates a slice with length 0 and room for `capacity` values. Using `make([]int, capacity)` instead would create `capacity` zeros that the stack would treat as real values.

### Memory after popping

Go code never frees memory directly. The **garbage collector** frees memory automatically once nothing in the program can reach it. The catch is that anything still reachable is never freed, even if the program will never use it again.

Shortening a slice does not free its underlying array. After `Pop`, the popped value is still stored in the array, just outside the slice's length, until a later `Push` overwrites it.

For a stack of `int`s this is harmless. For a stack that holds pointers or large values, the popped item can't be freed by the garbage collector while the array still refers to it. The usual fix is to clear the slot before shortening:

```go
top := s.items[len(s.items)-1]
s.items[len(s.items)-1] = nil // for a stack of pointers: release the reference
s.items = s.items[:len(s.items)-1]
```

The built-in `clear(s.items[len(s.items)-1:])` does the same, setting the slot to its type's zero value, and `slices.Delete` also zeroes the slots it removes.

The array also keeps its capacity after many pops. A stack that once held a million values keeps room for a million until the array is replaced. To give that memory back, copy the remaining values into a new, smaller slice, for example `s.items = slices.Clone(s.items)`. The old array then becomes unreachable, and the garbage collector frees it. Setting `s.items = nil` does the same when the stack is empty.

### Pitfalls

| Mistake | Concrete failure |
|---|---|
| Indexing without checking for an empty stack | `s.items[len(s.items)-1]` on an empty slice panics with `index out of range [-1]`. |
| Using a sentinel such as `-1` for "empty" | A stack that contains `-1` becomes indistinguishable from an empty one. |
| Using the front of the slice as the top | Every push and pop shifts all other values, so each operation costs O(n). |
| Shortening the slice before reading the top | The value is no longer inside the slice, so `Pop` can't return it. |
| `make([]int, n)` when only capacity was intended | The stack starts with `n` zeros on it. |

## Edge cases

| Category | Cases |
|---|---|
| Empty stack | `Pop` and `Peek` return `(0, false)` and change nothing; `Len` is 0 and `IsEmpty` is true. |
| One value | Popping it empties the stack, and a further `Pop` returns `(0, false)`. |
| Zero as a real value | Pushing `0` must give `(0, true)` from `Peek` and `Pop`, not the empty result. |
| `Peek` | Repeated calls return the same value and don't change `Len`. |
| Order | Pushing several values and popping them all returns them in reverse order. |
| Interleaving | Pushes and pops mixed together still follow LIFO. |
| Reuse | A stack that was emptied works normally when values are pushed again. |
| Growth | Thousands of pushes force `append` to grow the slice many times. |

Every size check compares `Len()` and `IsEmpty()`, so the two can't disagree.

## Variants

### Linked-list stack

A stack can also be built on a singly linked list, using the head as the top. Both ends of the stack's work happen at the front of the list, which is O(1):

```go
type ListStack struct {
    list linkedlist.LinkedList
}

func (s *ListStack) Push(v int) {
    s.list.AddFirst(v)
}

func (s *ListStack) Pop() (int, bool) {
    return s.list.DeleteFirst()
}
```

| | Slice-backed | Linked-list-backed |
|---|---|---|
| `Push` | amortized O(1) | O(1) every time |
| Memory per value | the value only | the value plus a pointer, and a separate allocation |
| Speed in practice | usually faster: values sit next to each other in memory | slower: each node is a separate allocation |
| Worst single push | O(n), when the slice grows | O(1) |

The slice version is the usual choice in Go. A linked list is worth it only when every single operation must be fast, with no occasional slow push.

### Generic stack

Go generics allow one stack type for any element type:

```go
type Stack[T any] struct {
    items []T
}
```

Methods then return `(T, bool)`, and the empty result is `T`'s zero value, which the boolean distinguishes from a real value, as before.

## When to use it

A stack fits any problem where the most recent item must be handled first:

- **Undo:** each action is pushed, and undo pops the latest one.
- **Matching brackets:** push each opening bracket, and pop when a closing bracket arrives. The popped bracket must match it.
- **Depth-first search:** an explicit stack replaces recursion, holding the nodes still to visit.
- **Function calls:** the program's call stack holds each active function's local variables, and returning pops them.

## Common mistakes

- **Calling every operation O(1).** `Push` is amortized O(1), because `append` occasionally copies the whole slice.
- **Picking a sentinel value for "empty".** Any `int` can be a real value. Return a boolean alongside the value.
- **Describing the front of the slice as "more complicated".** The precise reason to avoid it is cost: every operation at index 0 shifts all other values, which is O(n).

## Review questions

Try to answer each question before reading the answers in the next section.

1. What does LIFO mean, and which operations follow from it?
2. Why is the top of the stack the end of the slice, not the front?
3. Why is `Push` amortized O(1) rather than O(1), and what would make it O(n) on average?
4. How does `Pop` remove a value without copying, and why must it read the value first?
5. Why do `Pop` and `Peek` return `(int, bool)` rather than a special value?
6. Why does `&Stack{}` work without a constructor, and when would `New(capacity)` help?
7. What happens to a popped value's memory, and when does that matter?
8. How does a linked-list stack compare with a slice-backed one?
9. How would you use a stack to check whether brackets in a string are balanced?

## Answers

1. Last in, first out: the most recently pushed value is the first to be removed. That gives `Push` to add at the top, `Pop` to remove from the top, and `Peek` to read the top without removing it.
2. Adding or removing at the end of a slice moves nothing else, so it's O(1). At the front, every other value shifts by one position, which is O(n).
3. When the slice is full, `append` copies every value into a larger array, which is O(n) for that push. Capacity grows by a multiple, so those copies become rarer, and the total copying stays proportional to the number of pushes. Growing by a fixed amount would make the total copying grow with the square of the number of pushes, so the average push would be O(n).
4. It reslices to `s.items[:len(s.items)-1]`, a view that excludes the last element, without copying. After reslicing, the value is outside the slice, so it must be read first.
5. A stack of `int`s can hold any `int`, so no single value can mean "empty". The boolean separates "the top is 0" from "there is no top".
6. The `items` field starts as a `nil` slice, which has length 0 and works with `append`. `New(capacity)` reserves capacity up front, which avoids growth copies when the size is known in advance.
7. The value stays in the underlying array until it's overwritten, because shortening a slice doesn't free memory. It matters for pointers or large values, which can't be garbage collected while the array refers to them. Clearing the slot before shortening fixes it.
8. The linked-list version makes every push O(1), with no occasional slow push, but uses more memory per value and is usually slower in practice. The slice version is the default choice in Go.
9. Push each opening bracket. On each closing bracket, pop, and check that the popped bracket is the matching opening one. The string is balanced if every pop matches and the stack is empty at the end.
