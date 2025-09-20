package seq_test

import (
	"fmt"
	"iter"
	"slices"
	"testing"

	"github.com/sagikazarmark/seq"
)

func ExampleFlatten() {
	job1 := func() []string {
		return []string{"foo", "bar"}
	}
	job2 := func() []string {
		return []string{"baz"}
	}
	job3 := func() []string {
		return []string{"bat", "qux", "quux"}
	}

	doJobs := func() iter.Seq[[]string] {
		return slices.Values([][]string{job1(), job2(), job3()})
	}

	results := seq.Flatten(doJobs())

	for result := range results {
		fmt.Println(result)
	}

	// Output:
	// foo
	// bar
	// baz
	// bat
	// qux
	// quux
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
