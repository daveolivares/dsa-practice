# Doubly Linked List

Implementation: [`dlinkedlist.go`](dlinkedlist.go) · Tests: [`dlinkedlist_test.go`](dlinkedlist_test.go) · Builds on: [Singly Linked List](../linkedlist/LINKED_LIST.md)

## Core idea

A **doubly linked list** is a linked list where each node has two links: `next` to the node after it, and `prev` to the node before it. The list keeps a pointer to both ends: `head`, the first node, and `tail`, the last.

```text
        head                                tail
         │                                   │
         ▼                                   ▼
nil ◀──┤ 1 ├──▶◀──┤ 2 ├──▶◀──┤ 3 ├──▶ nil
       prev  next
```

The first node's `prev` and the last node's `next` are `nil`.

Two things change compared with a singly linked list:

- **Both ends are cheap.** Adding or removing at either end touches only a few pointers, so all four operations take O(1) time.
- **Any node can reach its neighbors directly.** A node you already hold can be removed in O(1), without searching for the node before it.

The cost is one extra pointer per node, and more pointers to keep consistent on every change.

## Invariants

These must be true before and after every method call:

1. **Links agree in both directions.** For every node `n`, if `n.next` is not `nil`, then `n.next.prev == n`. Stepping forward and then back returns to the same node.
2. **The ends are marked.** `head.prev == nil` and `tail.next == nil`.
3. **The ends agree on emptiness.** `head == nil`, `tail == nil`, and `length == 0` are either all true or all false.
4. **The count is accurate.** `length` equals the number of nodes.

In a one-element list, `head` and `tail` point to **the same node**. Many bugs come from code that assumes the two ends are different nodes.

Invariant 1 is the one most often broken, and it can be broken silently. If a method sets a `next` link but forgets the matching `prev` link, walking forward still looks correct, and only a backward walk reveals the problem.

## Complexity

`n` is the number of elements.

| Operation | Singly linked (no tail) | Doubly linked (head and tail) | Why it changed |
|---|---|---|---|
| `AddFirst`, `DeleteFirst` | O(1) | O(1) | |
| `AddLast` | O(n) | **O(1)** | The tail pointer gives direct access to the last node. |
| `DeleteLast` | O(n) | **O(1)** | `tail.prev` is the new last node, so no walk is needed. This requires both the tail pointer and `prev` links. |
| `Contains`, `DeleteValue` | O(n) | O(n) | Finding a value still means walking the list. |
| `Len` | O(1) | O(1) | |
| `ToSlice`, `ToSliceReverse`, `Reverse` | O(n) | O(n) | Each visits every node once. |

**Space:** O(n). Each element costs one node: the value plus two pointers.

A tail pointer alone does not make `DeleteLast` O(1) in a singly linked list, because the node before the tail can only be found by walking from the head. The `prev` link is what removes that walk.

## Implementation

### Types

```go
type node struct {
    value int
    prev  *node
    next  *node
}

type LinkedList struct {
    head   *node
    tail   *node
    length int
}
```

All fields are unexported, so only the list's own methods can change them and break the invariants. The zero value, `&LinkedList{}`, is a valid empty list.

### Adding at an end

Adding at the front of a non-empty list takes three pointer assignments:

1. the new node's `next` points to the old head,
2. the old head's `prev` points to the new node,
3. `head` moves to the new node.

```go
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
```

Step 1 happens when the node is created (`next: l.head`). On an empty list, `l.head` is `nil`, so the same line works, and the new node becomes the tail as well.

Step 2 is the one that is easy to forget, because singly linked lists don't need it. Without it, `[2, 3]` followed by `AddFirst(1)` walks forward as `[1, 2, 3]` but backward as `[3, 2]`.

The order matters: the old head must be updated **before** `head` moves, while `l.head` still refers to it.

`AddLast` is the mirror image, with `head` swapped for `tail` and `next` swapped for `prev`.

### Removing any node

Every removal uses one private helper. It unlinks a node `n` wherever `n` is in the list:

```go
func (l *LinkedList) remove(n *node) {
    if n.prev != nil {
        n.prev.next = n.next
    } else {
        l.head = n.next
    }

    if n.next != nil {
        n.next.prev = n.prev
    } else {
        l.tail = n.prev
    }

    // Detach n so the removed node no longer points into the list.
    n.prev = nil
    n.next = nil
    l.length--
}
```

The helper applies one rule to each side of `n`:

> **If `n` has a neighbor on this side, make the neighbor skip over `n`. If it doesn't, `n` is at this end of the list, so move the list's end pointer instead.**

The two sides are handled independently, which covers every position without special cases:

| Position of `n` | `n.prev` | `n.next` | Result |
|---|---|---|---|
| Middle | a node | a node | Both neighbors skip `n`. |
| Head | `nil` | a node | `head` moves forward; the new head's `prev` becomes `nil`. |
| Tail | a node | `nil` | `tail` moves back; the new tail's `next` becomes `nil`. |
| Only node | `nil` | `nil` | `head` and `tail` both become `nil`. |

Removing `2` from `[1, 2, 3]`:

```text
before:   [1] ⇄ [2] ⇄ [3]
                 n

            n.prev.next = n.next      (1 now points forward to 3)
            n.next.prev = n.prev      (3 now points back to 1)

after:    [1] ⇄ [3]
```

The three public delete methods then only decide **which** node to remove:

```go
func (l *LinkedList) DeleteFirst() (int, bool) {
    if l.head == nil {
        return 0, false
    }
    v := l.head.value
    l.remove(l.head)
    return v, true
}

func (l *LinkedList) DeleteValue(value int) bool {
    for curr := l.head; curr != nil; curr = curr.next {
        if curr.value == value {
            l.remove(curr)
            return true
        }
    }
    return false
}
```

`DeleteLast` is `DeleteFirst` with `l.tail` in place of `l.head`. The singly linked list needed a trailing `prev` variable in `DeleteValue`. Here every node already knows its neighbors, so a plain loop is enough.

### Working it out case by case

Before the helper existed, each delete method handled its cases explicitly. That version is a useful way to reason about the problem, even though the helper is the form to remember:

- **Empty list:** return `0, false`, and change nothing.
- **One element:** set `head` and `tail` to `nil`. Clearing only `head` leaves `tail` pointing at a removed node, and the next `AddLast` attaches to it.
- **Two or more elements, removing the head:** move `head` forward, then set the new head's `prev` to `nil`.
- **Middle node:** `curr.prev.next = curr.next` and `curr.next.prev = curr.prev`. These two lines are only safe when `curr` has a neighbor on both sides.

The helper's table above is exactly these cases, written as two independent checks.

### "First match" and duplicates

`DeleteValue` removes the first matching node, counting from the head. For adjacent duplicates the direction makes no visible difference: removing either `2` from `[1, 2, 2, 3]` leaves `[1, 2, 3]`. For **non-adjacent** duplicates it does: removing the first `2` from `[1, 2, 3, 2]` leaves `[1, 3, 2]`, while removing the last leaves `[1, 2, 3]`.

So "first match" is a real, testable requirement, and a test for it needs duplicates that are separated by other values. An earlier version of this code checked the tail before the middle nodes, which removed the last copy instead of the first. Tests with only adjacent duplicates did not catch it.

### Reversing in place

Reversing a doubly linked list swaps every node's `prev` and `next` links, then swaps `head` and `tail`:

```go
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
```

The step `curr = curr.prev` looks backwards but is correct. After the swap, `prev` holds what used to be `next`, so it still moves forward through the original order. Writing `curr = curr.next` instead follows the old `prev` link, which is `nil` at the head, so the loop stops after the first node. On `[1, 2, 3]`, that leaves a list that walks forward as `[3]`, and whose backward walk circles between two nodes forever.

Swapping `head` and `tail` unconditionally is safe. For an empty list they are both `nil`, and for one element they are the same node, so the swap changes nothing in either case.

The algorithm runs in O(n) time and O(1) extra space, and creates no nodes.

**Go-specific aside:** Go can swap two values in one statement, because the right-hand side is fully evaluated before anything is assigned. The swaps above could be written as `curr.next, curr.prev = curr.prev, curr.next` and `l.head, l.tail = l.tail, l.head`. The version with temporary variables is shown because it translates directly to any language.

### Testing the backward links

`ToSlice` only follows `next` links, so on its own it cannot detect a broken `prev` link. The tests also call `ToSliceReverse`, which walks from `tail` along `prev`, and check that it is exactly the reverse of `ToSlice`. A missing `prev` link, a stale `tail`, or a `head` whose `prev` is not `nil` all make the two walks disagree.

### Pitfalls

| Mistake | Concrete failure |
|---|---|
| Forgetting the old end's link when adding | `AddFirst(1)` on `[2, 3]` walks forward as `[1, 2, 3]` but backward as `[3, 2]`. |
| Moving `head` or `tail` before updating the old end node | In `AddFirst`, `l.head = newNode` followed by `l.head.prev = newNode` sets the new node's `prev` to itself, and the old head's `prev` is never updated. |
| Reading `head.next.prev` on a one-element list | `head.next` is `nil`, so the program crashes with a nil pointer dereference. |
| Clearing `head` but not `tail` when the list becomes empty | `tail` still points at the removed node, and the next `AddLast` attaches to it. |
| Checking the tail before the middle nodes in `DeleteValue` | `[1, 2, 3, 2]` loses its last `2` instead of its first. |
| A loop guard such as `curr.next != nil` without first checking `curr != nil` | Crashes when `curr` itself is `nil`, for example on a one-element list where the value isn't found. |
| `curr = curr.next` after swapping links in `Reverse` | The loop stops after the first node and leaves the list corrupted. |

## Edge cases

| Category | Why it matters |
|---|---|
| Empty list, for every operation | `head` and `tail` are `nil`, so any field access through them crashes. |
| One element | `head` and `tail` are the same node, so both must change together. |
| Two elements | The smallest list where the head and tail are different but adjacent. |
| Head, middle, and tail positions | Each position takes a different branch in `remove`. |
| Value not present, including in a one-element list | The search must reach `nil` and change nothing. |
| Adjacent and non-adjacent duplicates | Only non-adjacent duplicates can show which copy was removed. |
| Sequences of operations | Catch state left behind by earlier calls, such as a stale `tail` after emptying the list. |

Every test checks the forward walk, the backward walk, and `Len()`.

## Variants

### Sentinel node

A **sentinel** is a placeholder node that is always present and holds no real value. In a circular doubly linked list with a sentinel, the sentinel's `next` is the first real node and its `prev` is the last, so `head` and `tail` are just `sentinel.next` and `sentinel.prev`. Every real node always has two real neighbors, either of which may be the sentinel, so `remove` needs no `nil` checks. Go's `container/list` is built this way.

### Deque

A **deque** (double-ended queue) supports adding and removing at both ends. A doubly linked list provides all four operations in O(1), so it is a natural way to implement one. A slice-based ring buffer is the other common choice.

### LRU cache

An **LRU cache** (least recently used) keeps a fixed number of items and evicts the one that was used longest ago. It combines a hash map from key to node with a doubly linked list ordered by most recent use. Each access moves a node to the front, using `remove` followed by an add at the front, both O(1). Eviction removes the tail. This is one of the most common real-world uses of a doubly linked list.

## When to use it

- Use it when you need cheap adds and removes at **both ends**, or need to remove nodes you already hold pointers to, as in an LRU cache.
- If you only need one end, a singly linked list or a slice is simpler.
- For general sequences in Go, a **slice** is usually faster in practice, because its elements are stored next to each other in memory.

## Common mistakes

- **Treating `head` and `tail` as always different.** In a one-element list they are the same node.
- **Updating only one direction.** Every change to a `next` link needs the matching `prev` change, or the backward walk breaks while the forward walk still looks correct.
- **Relying on order of checks for correctness.** Special-casing the head and tail before searching the middle made the "first match" rule depend on the order of three separate checks. A single search from the head, followed by one shared `remove`, avoids that.
- **Forgetting the `prev` field itself.** A node with only `value` and `next` is a singly linked node.

## Review questions

Try to answer each question before reading the answers in the next section.

1. Which operations become O(1) compared with a singly linked list without a tail, and what makes each one possible?
2. State the invariant that relates `prev` and `next`. Why can breaking it go unnoticed?
3. Adding at the front of a non-empty list needs three pointer assignments. What are they, and in what order?
4. What does `remove(n)` do when `n` is the only node, and why does it need no special case?
5. Why does a one-element list make `head.next.prev = nil` crash?
6. When does "first match" in `DeleteValue` make a visible difference, and how would you test it?
7. In `Reverse`, why is the step `curr = curr.prev` rather than `curr = curr.next`?
8. Why is it safe for `Reverse` to swap `head` and `tail` without checking the list's size?
9. How can a test detect a broken `prev` link?
10. How does an LRU cache use a doubly linked list?

## Answers

1. `AddLast` (the tail pointer gives direct access to the last node) and `DeleteLast` (the tail pointer plus `tail.prev`, which is the new last node). `AddFirst` and `DeleteFirst` were already O(1).
2. For every node `n` with a next node, `n.next.prev == n`. Walking forward only follows `next`, so a missing or wrong `prev` link is invisible until something walks backward.
3. The new node's `next` points to the old head, the old head's `prev` points to the new node, then `head` moves to the new node. The old head must be updated before `head` moves.
4. `n.prev` and `n.next` are both `nil`, so the first check sets `head` to `nil` and the second sets `tail` to `nil`. Each side is handled independently, so the "only node" case falls out of the head and tail cases combined.
5. With one node, `head.next` is `nil`, and reading `.prev` through a `nil` pointer crashes. That case must set `head` and `tail` to `nil` instead.
6. Only when the duplicates are separated by other values, such as `[1, 2, 3, 2]`. Removing the first `2` leaves `[1, 3, 2]`, and removing the last leaves `[1, 2, 3]`. The test needs non-adjacent duplicates.
7. After swapping, `prev` holds the old `next`, so `curr.prev` is the next node in the original order. `curr.next` now holds the old `prev`, which leads backwards and is `nil` at the head.
8. For an empty list both are `nil`, and for one element both are the same node, so swapping them changes nothing.
9. Walk backward from `tail` along `prev` and check that the result is exactly the reverse of the forward walk.
10. A hash map finds a key's node in O(1). The list keeps nodes in order of recent use: an access removes the node and adds it at the front, and eviction removes the tail, all in O(1).
