package binarysearch

import "testing"

func TestSearch(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		target int
		want   int // -1 if absent; with duplicates, any matching index is accepted
	}{
		// Zero and one.
		{"empty", []int{}, 0, -1},
		{"single element found", []int{1}, 1, 0},
		{"single element not found", []int{1}, 2, -1},

		// Found at the boundaries and in the middle.
		{"first element", []int{1, 2, 3, 4, 5}, 1, 0},
		{"last element", []int{1, 2, 3, 4, 5}, 5, 4},
		{"middle element", []int{1, 2, 3, 4, 5}, 3, 2},
		{"even length", []int{1, 2, 3, 4, 5, 6}, 5, 4},

		// Not found, in every position it can be missing from.
		{"larger than all", []int{1, 2, 3, 4, 5}, 6, -1},
		{"smaller than all", []int{1, 2, 3, 4, 5}, 0, -1},
		{"between elements", []int{1, 3, 5, 7}, 4, -1},

		// Duplicates: any matching index is correct.
		{"duplicates", []int{1, 2, 2, 3, 4}, 2, 1},
		{"all duplicates", []int{2, 2, 2, 2, 2}, 2, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Search(tt.nums, tt.target)

			if tt.want == -1 {
				if got != -1 {
					t.Errorf("Search(%v, %d) = %d, want -1", tt.nums, tt.target, got)
				}
				return
			}

			// Check the requirement, not one implementation's choice of index.
			if got < 0 || got >= len(tt.nums) || tt.nums[got] != tt.target {
				t.Errorf("Search(%v, %d) = %d, want an index i with nums[i] == %d",
					tt.nums, tt.target, got, tt.target)
			}
		})
	}
}
