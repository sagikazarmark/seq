package seq_test

import (
	"fmt"
	"iter"
	"slices"
	"testing"

	"github.com/sagikazarmark/seq"
)

func ExampleFlatten() {
	s := slices.Values([][]string{
		{"hello", "world"},
		{"how", "are", "you"},
		{"today"},
	})

	greeting := seq.Flatten(s)

	for word := range greeting {
		fmt.Println(word)
	}

	// Output:
	// hello
	// world
	// how
	// are
	// you
	// today
}

func TestFlatten(t *testing.T) {
	testCases := []struct {
		name     string
		seq      iter.Seq[[]int]
		expected []int
	}{
		{"empty_sequence", slices.Values([][]int{}), []int{}},
		{"single_empty_slice", slices.Values([][]int{{}}), []int{}},
		{"single_slice_with_elements", slices.Values([][]int{{1, 2, 3}}), []int{1, 2, 3}},
		{"multiple_slices", slices.Values([][]int{{1, 2}, {3, 4}, {5}}), []int{1, 2, 3, 4, 5}},
		{"empty_slices_mixed_with_non_empty", slices.Values([][]int{{1, 2}, {}, {3}}), []int{1, 2, 3}},
		{"all_empty_slices", slices.Values([][]int{{}, {}, {}}), []int{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := slices.Collect(seq.Flatten(tc.seq))
			if !slices.Equal(result, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}
