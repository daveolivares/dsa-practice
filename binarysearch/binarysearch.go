// Package binarysearch implements binary search over sorted slices.
package binarysearch

// Search returns the index of target in nums, or -1 if it is not present.
// nums must be sorted in ascending order. If target appears more than once,
// any matching index may be returned.
//
// Time: O(log n). Space: O(1).
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

// SearchFirst returns the index of the first occurrence of target in nums,
// or -1 if it is not present. nums must be sorted in ascending order.
//
// Time: O(log n). Space: O(1).
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
