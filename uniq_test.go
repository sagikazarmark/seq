package seq_test

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"testing"

	"github.com/sagikazarmark/seq"
)

func ExampleUniq() {
	numbers := slices.Values([]int{1, 2, 2, 3, 1, 4, 3, 5})

	unique := seq.Uniq(numbers)

	for number := range unique {
		fmt.Println(number)
	}

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

func TestUniq(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"empty_sequence", []int{}, []int{}},
		{"no_duplicates", []int{1, 2, 3, 4}, []int{1, 2, 3, 4}},
		{"all_duplicates", []int{1, 1, 1, 1}, []int{1}},
		{"mixed_duplicates", []int{1, 2, 2, 3, 1, 4, 3, 5}, []int{1, 2, 3, 4, 5}},
		{"consecutive_duplicates", []int{1, 1, 2, 2, 3, 3}, []int{1, 2, 3}},
		{"single_element", []int{42}, []int{42}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := slices.Values(tc.input)
			actual := slices.Collect(seq.Uniq(input))

			if !slices.Equal(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func ExampleUniq2() {
	roles := seq.Chain2(
		maps.All(map[string]string{"alice": "admin", "bob": "user"}),
		maps.All(map[string]string{"alice": "user", "charlie": "user"}),
	)

	unique := seq.Uniq2(roles)

	printSorted(unique)

	// Output:
	// alice: admin
	// bob: user
	// charlie: user
}

func TestUniq2(t *testing.T) {
	testCases := []struct {
		name     string
		input    iter.Seq2[string, int]
		expected map[string]int
	}{
		{
			"empty_sequence",
			maps.All(map[string]int{}),
			map[string]int{},
		},
		{
			"no_duplicate_keys",
			maps.All(map[string]int{"a": 1, "b": 2, "c": 3}),
			map[string]int{"a": 1, "b": 2, "c": 3},
		},
		{
			"single_pair",
			maps.All(map[string]int{"key": 42}),
			map[string]int{"key": 42},
		},
		{
			"duplicate_keys_from_manual_seq",
			seq.Chain2(
				maps.All(map[string]int{"a": 1, "b": 2}),
				maps.All(map[string]int{"a": 3, "c": 4, "b": 5}),
			),
			map[string]int{"a": 1, "b": 2, "c": 4},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := maps.Collect(seq.Uniq2(tc.input))

			if !maps.Equal(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}
