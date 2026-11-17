# Singly Linked List

Implementation: [`linkedlist.go`](linkedlist.go) · Tests: [`linkedlist_test.go`](linkedlist_test.go)

## Core idea

A **singly linked list** stores a sequence of values in separate **nodes**. Each node holds one value and a pointer to the next node. The list itself only needs to know where the first node is, called the **head**. The last node's `next` pointer is `nil`, which marks the end.

```text
head
 │
 ▼
┌───┬───┐    ┌───┬───┐    ┌───┬─────┐
│ 1 │ ●─┼───▶│ 2 │ ●─┼───▶│ 3 │ nil │
└───┴───┘    └───┴───┘    └───┴─────┘
```

A slice stores its elements side by side in memory, so it can jump straight to any index. A linked list's nodes can be anywhere in memory, and the only way to reach a node is to start at the head and follow `next` pointers one at a time. That trade-off explains every cost in this structure:

- **Changing the front is cheap.** Adding or removing the head only touches one or two pointers, whatever the list's length.
- **Reaching anything else is expensive.** Finding the last node, a specific value, or the node at position `i` means walking the list.

"Singly" means each node points in one direction only, forward. There is no way to step back from a node to the one before it.

## Invariants

These statements must be true before and after every method call. Each method is correct only if it keeps all of them true:

1. **The list ends.** Following `next` from `head` reaches `nil` after a finite number of steps. There are no cycles, meaning no node points back to an earlier node or to itself.
2. **The count is accurate.** `length` equals the number of nodes reachable from `head`.
3. **Empty means `head == nil`.** A list with no nodes has `head == nil` and `length == 0`, and vice versa.

Most linked-list bugs break one of these. A node that points to itself breaks invariant 1, and every walk of the list then runs forever. A counter update that a code path skips breaks invariant 2.

## Complexity

`n` is the number of elements.

| Operation | Time | Why |
|---|---|---|
| `AddFirst` | O(1) | Creates one node and moves `head`. |
| `DeleteFirst` | O(1) | Moves `head` to the second node. |
| `Len` | O(1) | Reads the stored `length` field instead of counting nodes. |
| `AddLast` | O(n) | Must walk to the last node, because there is no tail pointer. |
| `DeleteLast` | O(n) | Must walk to the node *before* the last one, to set its `next` to `nil`. |
| `Contains` | O(n) | May check every node. It stops early if it finds the value. |
| `DeleteValue` | O(n) | May check every node to find the first match. |
| `ToSlice` | O(n) | Visits every node once. |

**Space:** O(n). Each element costs one node, which holds the value plus one pointer.

## Implementation

### Types

```go
type node struct {
    value int
    next  *node
}

type LinkedList struct {
    head   *node
    length int
}
```

All of these names start with a lowercase letter, so they are **unexported**: code outside the package cannot see them. This protects the invariants. If callers could set `length` or `head` directly, they could put the list into an impossible state. They use the methods instead.

The **zero value** works without a constructor. `&LinkedList{}` has `head == nil` and `length == 0`, which is a valid empty list.

### Walking the list

Almost every operation is built on the same loop. It starts at `head` and follows `next` until it reaches `nil`:

```go
curr := l.head
for curr != nil {
    // use curr.value
    curr = curr.next
}
```

On an empty list, `curr` starts as `nil` and the loop body never runs, so `Contains` and `ToSlice` need no special case for an empty list.

### Changing the front: `AddFirst` and `DeleteFirst`

```go
func (l *LinkedList) AddFirst(value int) {
    newNode := &node{value: value, next: l.head}
    l.head = newNode
    l.length++
}
```

The order matters. The new node must capture the **old** head before `head` moves. On an empty list, `l.head` is `nil`, so the new node's `next` is `nil` and the same code handles both cases.

`DeleteFirst` moves `head` to `head.next`. When the list has one element, `head.next` is `nil`, so the list correctly becomes empty.

### Changing the end: `AddLast` and `DeleteLast`

`AddLast` walks to the last node, the one whose `next` is `nil`, and attaches the new node there. It needs a special case for an empty list, because there is no last node to attach to:

```go
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
```

`DeleteLast` must stop one node earlier, at the **second-to-last** node, because that is the node whose `next` has to become `nil`. The loop condition `curr.next.next != nil` stops there. A one-element list has no second-to-last node, so it is handled before the loop:

```go
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
```

The one-element check must come before the loop. With a single node, `curr.next` is `nil`, and evaluating `curr.next.next` would crash.

### Removing from the middle: `DeleteValue`

To remove a node, the node **before** it must skip over it. A singly linked list cannot step backward, so the loop keeps two pointers: `curr` is the node being checked, and `prev` stays one node behind.

```go
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
```

Removing `2` from `[1, 2, 3]`:

```text
before:   prev        curr
           │           │
           ▼           ▼
head ──▶ [1|●] ───▶ [2|●] ───▶ [3|nil]

after:    prev.next = curr.next
head ──▶ [1|●] ──────────────▶ [3|nil]
```

Node `2` is no longer reachable from `head`, so it is no longer part of the list, and Go's garbage collector frees it.

The head needs a separate case because it has no previous node. Removing it means moving `head` instead of changing a `next` field. Reusing `DeleteFirst` keeps that logic in one place. Returning straight after the first removal implements the "first match only" rule: later duplicates are never reached.

### Updating the counter

`length` must change exactly when a node is added or removed, and only then:

- Every successful add increments it once, on every path through the method.
- Every delete decrements it only on paths that actually removed a node. The empty-list and value-not-found paths must leave it unchanged.

### Pitfalls

| Mistake | Concrete failure |
|---|---|
| Walking from `head` without checking for an empty list | In `AddLast`, `curr := l.head` followed by `curr.next` crashes on an empty list with "nil pointer dereference", because `curr` is `nil`. |
| Handling the empty case without stopping | If the empty-list branch of `AddLast` sets `head = newNode` and then continues into the general path, it finds `curr.next == nil` on the new node and sets `curr.next = newNode`. The node now points to itself, and `ToSlice` loops forever. |
| An early `return` that skips the counter update | `AddLast` adds the node but `Len()` reports one less than the real count. |
| Decrementing before checking that a node exists | A delete on an empty list makes `length` negative. |
| Assigning in the wrong order | In `AddFirst`, `l.head = newNode` followed by `newNode.next = l.head` makes the new node point to itself, and the old list is lost. |
| `curr.next.next` on a one-element list | Crashes, because `curr.next` is `nil`. |
| Forgetting the head case in `DeleteValue` | Starting the `prev`/`curr` loop at the second node means the head can never be removed. |

When a test run hangs instead of failing, suspect a cycle. `go test -timeout 5s` makes it fail quickly with a stack trace instead.

## Edge cases

The tests are grouped by these categories, which apply to most linked structures:

| Category | Why it matters |
|---|---|
| Empty list, for every operation | The `head == nil` case, where any `curr.next` crashes. |
| One element | The head is also the last node, so "remove the head" and "remove the last" are the same change. |
| Two elements | The smallest list with a second-to-last node, which `DeleteLast`'s loop depends on. |
| Head, middle, and last positions | Each one uses a different code path in `DeleteValue`. |
| Value not present | The loop must reach `nil` and change nothing. |
| Duplicates | Checks the "first match only" rule. |
| Sequences of operations | Catch state left behind by earlier calls, such as reusing a list after emptying it. |

Every test checks `Len()` together with the contents, so a counter that drifts out of sync fails immediately, whichever operation caused it.

## Variants

### Tail pointer

Adding a `tail` field that points to the last node makes `AddLast` O(1).

It does **not** make `DeleteLast` O(1). After removing the last node, the new last node is the one before it, and a singly linked list can only find that node by walking from `head`.

The cost is extra bookkeeping. Every operation that can change which node is last must update `tail`: adding to an empty list, removing the only element, `DeleteLast`, and `DeleteValue` when it removes the last node. Forgetting one leaves `tail` pointing at a node that is no longer in the list.

A singly linked list with a tail pointer is a natural **queue**: add at the back and remove from the front, both in O(1).

### Doubly linked list

Each node also stores `prev`, a pointer to the node before it. With head and tail pointers, both ends support O(1) adds and removes, and removing a node you already hold a pointer to is O(1), because its neighbors are reachable directly.

The costs are one extra pointer per node and twice as many pointers to keep consistent on every change. Go's standard library `container/list` is a doubly linked list.

### Removing without a head special case

`DeleteValue` treats the head separately because removing the head changes `l.head`, while removing any other node changes the `next` field of the node before it. Two common techniques remove that difference. Both produce the same results as the `prev`/`curr` version.

#### Dummy node

A **dummy node**, also called a **sentinel**, is a placeholder node placed in front of the real head. It holds no real value. Because the dummy sits before the head, every real node, including the head, now has a node before it, so one loop handles every position:

```go
func (l *LinkedList) DeleteValue(value int) bool {
    dummy := &node{next: l.head} // placeholder in front of the real head
    prev := dummy
    for prev.next != nil {
        if prev.next.value == value {
            prev.next = prev.next.next
            l.head = dummy.next // the head may have been the node removed
            l.length--
            return true
        }
        prev = prev.next
    }
    return false
}
```

Two details matter here:

- The loop checks `prev.next`, the node *after* `prev`, so `prev` is always the node whose `next` needs to change.
- After a removal, `l.head` is reset from `dummy.next`. If the head was removed, `dummy.next` is now the second node, and the list's `head` must follow it.

This version creates a temporary dummy node on each call. Some implementations instead store a **permanent** sentinel in the list, so `head` is always the dummy and the real elements start at `head.next`. That removes head special cases from every method, at the cost of one extra node and code that has to skip over it. Go's `container/list` uses a permanent sentinel.

#### Pointer to the incoming pointer

Every node is reached through exactly one pointer: either `l.head` or the `next` field of the node before it. This version tracks the **address of that incoming pointer** instead of tracking the previous node:

```go
func (l *LinkedList) DeleteValue(value int) bool {
    incoming := &l.head // address of the pointer that leads to the current node
    for *incoming != nil {
        if (*incoming).value == value {
            *incoming = (*incoming).next // redirect that pointer past the node
            l.length--
            return true
        }
        incoming = &(*incoming).next
    }
    return false
}
```

`incoming` has type `**node`, a pointer to a pointer. It starts as the address of `l.head`, then becomes the address of each node's `next` field in turn. `*incoming` is the current node. Assigning to `*incoming` changes whichever pointer led to that node, whether it is `l.head` or some node's `next`, so the head needs no special treatment and no extra node is created.

#### Which to use

| Version | Head special case | Extra node | Readability |
|---|---|---|---|
| `prev`/`curr` (this implementation) | Yes | No | Most familiar |
| Dummy node | No | One per call, or one permanent | Familiar, slightly more setup |
| Pointer to the incoming pointer | No | No | Shortest, but double pointers are unusual in Go |

All three are correct and run in O(n). The `prev`/`curr` version is the most common in Go code. The dummy node is a standard technique in interview problems and in code that edits lists in many places.

## When to use it

In Go, a **slice is usually the better default**. Its elements sit next to each other in memory, which makes walking through them much faster in practice. `append` also adds to the end in amortized O(1) time, meaning it is O(1) on average across many calls.

A linked list is the better fit when:

- you add or remove at the **front** often, which costs O(n) in a slice, because every element has to shift,
- you insert or remove next to nodes you **already hold pointers to**, using a doubly linked list,
- you are building another structure on top of it, such as a stack, a queue, an LRU cache (doubly linked list plus hash map), or the chains in a hash table.

## Common mistakes

- **Assuming a tail pointer makes `DeleteLast` O(1).** It only speeds up adding at the end. Removing the last node still requires finding the node before it.
- **Treating the empty list as an afterthought.** Every method that touches `curr.next` needs a plan for `head == nil`. In practice, that is the first case to write.
- **Writing a special case that doesn't stop.** Handling the empty list in `AddLast` and then falling through into the general path created a node that pointed to itself.
- **Updating a counter on only some paths.** An early `return` skipped `length++`. Put the update where every successful path reaches it, and nowhere that a failed operation reaches.

## Review questions

Try to answer each question before reading the answers in the next section.

1. Why is `AddFirst` O(1) but `AddLast` O(n) in this implementation?
2. What are the list's three invariants, and which bug breaks each one?
3. In `AddFirst`, why must the new node's `next` be set before `head` changes?
4. Why does `DeleteLast` stop at the second-to-last node, and what loop condition does that?
5. Why does a one-element list need its own case in `DeleteLast`?
6. Why does `DeleteValue` need a `prev` pointer, and why is the head handled separately?
7. With a tail pointer, which operations become faster, and which do not?
8. A test run hangs instead of failing. What is the likely cause, and how do you make it fail fast?
9. When would you choose a slice over a linked list in Go?
10. How does a dummy node remove the head special case in `DeleteValue`, and what must happen to `l.head` afterwards?

## Answers

1. `AddFirst` only changes `head` and the new node. `AddLast` has to walk every node to find the last one, because the list stores no pointer to it.
2. The list ends (no cycles), `length` matches the number of nodes, and empty means `head == nil`. A node pointing to itself breaks the first. A skipped counter update breaks the second. Setting `head` to `nil` without resetting `length`, or the reverse, breaks the third.
3. The new node must point to the old first node. If `head` is moved first, the old first node is no longer reachable, and the new node ends up pointing to itself.
4. The node whose `next` must become `nil` is the second-to-last one. `curr.next.next != nil` keeps moving while there are at least two nodes ahead, so it stops with `curr.next` being the last node.
5. A single node has no node before it, so the only change is setting `head` to `nil`. Running the loop instead would evaluate `curr.next.next` with `curr.next` equal to `nil`, which crashes.
6. Removing a node means changing the `next` of the node before it, and a singly linked list cannot step backward, so `prev` remembers that node. The head has no node before it, so removing it changes `head` instead.
7. `AddLast` becomes O(1). `DeleteLast` stays O(n), because it still has to find the node before the last one. `AddFirst`, `DeleteFirst`, and `Len` were already O(1).
8. A cycle, usually a node whose `next` points to itself or to an earlier node, so walking the list never reaches `nil`. Run `go test -timeout 5s` so it fails with a stack trace.
9. Usually by default: slices are faster to walk in practice and `append` is amortized O(1). A linked list wins when you frequently add or remove at the front or next to nodes you already hold.
10. The dummy node sits in front of the head, so the head has a node before it like every other node, and one loop handles all positions. Afterwards `l.head` must be set to `dummy.next`, in case the head was the node removed.
