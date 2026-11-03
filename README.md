# dsa-practice

Data structures and algorithms implemented from scratch in Go, for understanding and long-term retention.

## How this repo works

- **One package per topic** (for example `binarysearch/`, `linkedlist/`, `stack/`), all in a single Go module.
- **Spec and tests first.** Each topic starts by deciding its operations, edge cases, and error behavior, then writing table-driven tests, then implementing.
- Each package's doc comment states the operations and their time and space complexity.
- **Study notes** (`<TOPIC>.md`, for example `binarysearch/BINARY_SEARCH.md`) are written after a topic is implemented and reviewed, for later reference.

## Running tests

```sh
go test ./...                # every topic
go test ./binarysearch/...   # one topic
```

## Topics

### Searching

- [x] Binary search

### Linear structures (review)

- [ ] Singly linked list
- [ ] Doubly linked list
- [ ] Stack
- [ ] Queue (linked and circular buffer)

### Hashing (review)

- [ ] Hash table
- [ ] Hash set

### Trees (review)

- [ ] Binary tree and traversals
- [ ] Binary search tree
- [ ] Trie

### New material

- [ ] Binary heap and priority queue
- [ ] Graphs: representations, BFS/DFS, topological sort, union-find, Dijkstra
- [ ] Sorting: merge sort, quicksort, heap sort
- [ ] Patterns: two pointers, sliding window, fast/slow pointers, backtracking, greedy, intervals
- [ ] Dynamic programming
