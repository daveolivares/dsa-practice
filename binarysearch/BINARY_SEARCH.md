# Binary Search

Implementation: [`binarysearch.go`](binarysearch.go) · Tests: [`binarysearch_test.go`](binarysearch_test.go)

## Core idea

Binary search finds a target value in a **sorted** slice. It keeps track of a **search range**, the part of the slice where the target could still be, using two indices, `low` and `high`. On each step it:

1. looks at the element in the middle of the range,
2. stops if that element is the target,
3. otherwise discards the half of the range that cannot contain the target.

Discarding is safe only because the slice is sorted. If the middle element is smaller than the target, every element to its left is smaller too.

The slice itself is never changed or copied. Only `low` and `high` move.

### Example: target found

Searching for `9` in `[1, 3, 5, 7, 9, 11, 13]` (indices 0–6):

| Step | `low` | `high` | `middle` | `nums[middle]` | Comparison | Next |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | 0 | 6 | 3 | 7 | 7 < 9 | `low = 4` |
| 2 | 4 | 6 | 5 | 11 | 11 > 9 | `high = 4` |
| 3 | 4 | 4 | 4 | 9 | 9 == 9 | return 4 |

### Example: target missing

Searching for `8` in the same slice:

| Step | `low` | `high` | `middle` | `nums[middle]` | Comparison | Next |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | 0 | 6 | 3 | 7 | 7 < 8 | `low = 4` |
| 2 | 4 | 6 | 5 | 11 | 11 > 8 | `high = 4` |
| 3 | 4 | 4 | 4 | 9 | 9 > 8 | `high = 3` |
| 4 | 4 | 3 | — | — | `low > high` | return -1 |

The range is empty after step 3, so the loop stops.

## Invariant

An **invariant** is a statement that is true at the start of every loop iteration. For binary search:

> **If the target is in `nums`, its index is between `low` and `high`, inclusive.**

This is what makes the algorithm correct:

- **It is true before the first iteration**, because the range `[0, len(nums)-1]` is the whole slice.
- **Every update keeps it true.** If `nums[middle] < target`, then `nums[middle]` and everything to its left are too small, so `low = middle + 1` removes only positions that cannot hold the target. The same reasoning applies to `high = middle - 1` when `nums[middle] > target`.
- **When the loop ends, the range is empty** (`low > high`). The target must be inside an empty range, which is impossible, so it is not in the slice and returning -1 is correct.

Do not confuse the invariant with the **loop condition**, `low <= high`. The condition decides whether to keep going. The invariant explains why the result is right when the loop stops.

## Complexity

- **Time: O(log n).** Each iteration removes at least half of the remaining range.
- **Space: O(1).** It uses a fixed number of variables, regardless of input size.

`log₂ n` is the number of times `n` can be halved before reaching 1. For example, 16 → 8 → 4 → 2 → 1 is four halvings, so `log₂ 16 = 4`. The loop runs at most `⌊log₂ n⌋ + 1` times:

| Slice length       | 16 | 1,000 | 1,000,000 | 1,000,000,000 |
| ------------------ | -- | ----- | --------- | ------------- |
| Maximum iterations | 5  | 10    | 20        | 30            |

Doubling the slice length adds only one iteration. A linear search (checking every element) takes up to `n` iterations: one million for one million elements.

## Implementation

This version uses an **inclusive** range, meaning `high` is the index of the last candidate, not one past it:

```go
func Search(nums []int, target int) int {
    low, high := 0, len(nums)-1
    for low <= high {
        middle := low + (high-low)/2

        if nums[middle] == target {
            return middle
        }

        if nums[middle] < target {
            low = middle + 1
        } else {
            high = middle - 1
        }
    }

    return -1
}
```

Three decisions must agree with one another:

| Decision | Choice here | Why |
| --- | --- | --- |
| Initial range | `high = len(nums)-1` | `high` points at the last element, so the range includes both ends. |
| Loop condition | `low <= high` | When `low == high`, one element remains and must still be checked. |
| Updates | `middle + 1` and `middle - 1` | `middle` has just been checked, so it can be excluded. This also guarantees the range shrinks on every iteration. |

An empty slice needs no special handling: `high` starts at -1, so the loop body never runs.

### Pitfalls

| Mistake | Concrete failure |
| --- | --- |
| Loop condition `low < high` with an inclusive range | For `[5]` with target 5, `low` and `high` are both 0, so the loop never runs and the function returns -1. |
| Updating with `high = middle` while using `low <= high` | For `[1, 3]` with target 0, the range reaches `low = high = 0`, and `high = middle` sets `high` to 0 again. The range never shrinks, so the loop runs forever. |
| Computing `middle := (low + high) / 2` | `low + high` can exceed the largest value the integer type holds and wrap around to a negative number (**integer overflow**). With 32-bit integers, `low = high = 1,500,000,000` overflows. Java's standard library shipped this bug for about nine years. Go's `int` is 64-bit on modern machines, so it is very unlikely there, but `low + (high-low)/2` is safe everywhere. |
| Mixing inclusive and half-open conventions | Each line can look correct on its own while the combination skips elements or loops forever. Choose one convention and apply it to the initial range, the loop condition, and every update. The half-open convention appears under [Lower bound](#lower-bound-insertion-point). |

## Edge cases

The tests are grouped by these categories, which apply to most search problems:

| Category | Cases |
| --- | --- |
| Zero and one | Empty slice; single element, found and not found |
| Position of the target | First index, last index, middle |
| Slice length | Odd and even lengths, which place the middle differently |
| Missing target | Smaller than every element, larger than every element, between two elements. Each reaches an empty range by a different path. |
| Duplicates | In the middle, at the start, at the end, filling the whole slice, and duplicates of a value other than the target |

### Testing when several answers are correct

When the target appears more than once, `Search` may return the index of any copy. Which copy it returns depends on details such as how `middle` is rounded. Two correct implementations can return different indices for the same input, each one consistently.

A test should therefore check the requirement, not one implementation's result. Instead of `got == 1`, the test checks that `got` is a valid index and `nums[got] == target`. For `SearchFirst`, which has exactly one correct answer, the test compares indices directly.

## Variants

The variants below differ from `Search` mainly in one decision: **when `nums[middle] == target`, which way should the range move?**

### First occurrence (`SearchFirst`)

On a match, record the index and keep searching **left**, because an earlier copy can only be on the left.

```go
func SearchFirst(nums []int, target int) int {
    low, high := 0, len(nums)-1
    result := -1
    for low <= high {
        middle := low + (high-low)/2

        if nums[middle] == target {
            result = middle
        }

        if nums[middle] < target {
            low = middle + 1
        } else {
            // Also runs on a match: keep searching left for an earlier one.
            high = middle - 1
        }
    }

    return result
}
```

On a match, `nums[middle] < target` is false, so the `else` branch moves left. Each match recorded is further left than the previous one, so the last value of `result` is the first occurrence.

Searching for `2` in `[1, 2, 2, 2, 3]`:

| Step | `low` | `high` | `middle` | `nums[middle]` | `result` | Next |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | 0 | 4 | 2 | 2 (match) | 2 | `high = 1` |
| 2 | 0 | 1 | 0 | 1 | 2 | `low = 1` |
| 3 | 1 | 1 | 1 | 2 (match) | 1 | `high = 0` |
| 4 | 1 | 0 | | | 1 | return 1 |

The range still halves on every iteration, including matches, so the cost stays **O(log n)** even when every element equals the target.

### Last occurrence

This is the mirror image: on a match, keep searching **right**. The only change from `SearchFirst` is `<` becoming `<=`, which sends a match into the branch that moves right:

```go
if nums[middle] == target {
    result = middle
}

if nums[middle] <= target {
    low = middle + 1 // also runs on a match: keep searching right
} else {
    high = middle - 1
}
```

With the same input, the matches are recorded at indices 2 and then 3, and the function returns 3.

First and last occurrence together give the number of copies of a value in O(log n): `last - first + 1`.

### Lower bound (insertion point)

This variant returns the first index `i` where `nums[i] >= target`. If the target is present, that is its first occurrence. If not, it is the index where the target would be inserted to keep the slice sorted, which can be `len(nums)`.

It is usually written with a **half-open** range, `[low, high)`, where `low` is included but `high` is one past the last candidate:

```go
func LowerBound(nums []int, target int) int {
    low, high := 0, len(nums)
    for low < high {
        middle := low + (high-low)/2

        if nums[middle] < target {
            low = middle + 1
        } else {
            high = middle
        }
    }

    return low
}
```

Each of the three decisions changes to fit the half-open range. `high` starts at `len(nums)`. The loop condition is `low < high`, because `low == high` means the range is empty. The update is `high = middle`, not `middle - 1`, because `nums[middle] >= target` means `middle` might be the answer and must remain a possible result.

Searching for `4` in `[1, 3, 5, 7]`:

| Step | `low` | `high` | `middle` | `nums[middle]` | Next |
| --- | --- | --- | --- | --- | --- |
| 1 | 0 | 4 | 2 | 5 (≥ 4) | `high = 2` |
| 2 | 0 | 2 | 1 | 3 (< 4) | `low = 2` |
| 3 | 2 | 2 | | | return 2 |

Index 2 is between 3 and 5, where 4 would be inserted.

Go's standard library provides this as `slices.BinarySearch(nums, target)`, which returns the position and whether the target was found. The older `sort.Search` does the same for any condition.

### Searching over answers instead of a slice

Binary search can also search a range of numbers for the boundary where a condition changes from true to false. This works when the condition is **monotonic**, meaning that once it becomes false, it stays false for every larger value.

For example, the integer square root of `n` is the largest `x` where `x*x <= n`. The condition is true for small `x` and false after the boundary, so it has the same shape as last occurrence:

```go
// IntSqrt returns the largest x with x*x <= n, for n >= 0.
// It assumes n is small enough that middle*middle does not overflow.
func IntSqrt(n int) int {
    low, high := 0, n
    result := 0
    for low <= high {
        middle := low + (high-low)/2

        if middle*middle <= n {
            result = middle
            low = middle + 1
        } else {
            high = middle - 1
        }
    }

    return result
}
```

The same pattern answers questions like "what is the smallest capacity that finishes the job in time?" Binary search is useful whenever you can check a candidate answer cheaply and the answers form a true/false boundary.

## When to use it

- **The data is sorted, or will be searched many times.** Sorting costs O(n log n), more than a single O(n) linear search, so sorting just to search once is not worth it.
- **You need ordered queries.** A hash map finds an exact key in O(1) on average, but cannot answer "what is the first value ≥ x?" or "how many values fall between a and b?" A sorted slice with binary search answers both in O(log n).
- **The problem has a true/false boundary.** Phrases such as "sorted", "first or last position", or "smallest value that satisfies" often point to binary search.

## Common mistakes

- **Giving the complexity as O(n log n).** That is the typical cost of *sorting*, and it is worse than linear search. Binary search is O(log n).
- **Assuming the basic version returns the first occurrence.** It returns whichever copy it reaches first. Finding the first occurrence requires continuing to search after a match.
- **Calling the loop condition the invariant.** `low <= high` controls the loop. The invariant, "if present, the target is in `[low, high]`", is what proves the answer correct.
- **Relying on code that works by accident.** In `SearchFirst`, the match case is handled by falling into the `else` branch. If that is not deliberate and commented, a later edit, such as adding `else if nums[middle] > target`, can break it quietly. Tests for duplicates catch that kind of change.

## Review questions

Try to answer each question before reading the answers in the next section.

1. Why must the slice be sorted?
2. What is the invariant, and why does it make returning -1 correct?
3. What is the maximum number of iterations for 1,000,000 elements? For 2,000,000?
4. With an inclusive range, what goes wrong if the loop condition is `low < high`?
5. Why compute the middle as `low + (high-low)/2`?
6. In `SearchFirst`, why move left on a match, and why is the last recorded match the first occurrence?
7. What single change turns `SearchFirst` into a last-occurrence search?
8. In the half-open lower bound, why is the update `high = middle` instead of `middle - 1`?
9. How do you test a function that may return any of several correct answers?

## Answers

1. Sorting is what makes discarding half the range safe. If `nums[middle] < target`, every element to the left of `middle` is also smaller than the target, so none of them can be the target.
2. If the target is in the slice, its index is in `[low, high]`. The loop ends only when `low > high`, so the range is empty, and the target cannot be in the slice.
3. 20 and 21. Doubling the input adds one halving.
4. The loop stops when one element remains, without checking it. For example, `[5]` with target 5 returns -1.
5. `low + high` can overflow the integer type and become negative. `high - low` is never larger than the slice length, so it cannot overflow.
6. Every element to the right of a match is greater than or equal to the target, so an earlier copy can only be on the left. Each new match is further left than the previous one, so the last match recorded is the leftmost.
7. Change `nums[middle] < target` to `nums[middle] <= target`, so a match moves the range right instead of left.
8. The answer is always between `low` and `high`, with `high` itself included, because the answer can be `len(nums)`. `nums[middle] >= target` means `middle` could be the answer. `high = middle` keeps it as a possible answer; `high = middle - 1` could remove the correct answer. The range still shrinks because `middle` is always less than `high`.
9. Check a property every correct answer has, such as `nums[got] == target`, rather than comparing against one specific value.
