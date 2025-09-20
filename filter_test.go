package seq_test

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"testing"

	"github.com/sagikazarmark/seq"
)

func ExampleFilter() {
	numbers := slices.Values([]int{1, 2, 3, 4, 5})

	odd := func(n int) bool {
		return n%2 == 1
	}

	oddNumbers := seq.Filter(numbers, odd)

	for n := range oddNumbers {
		fmt.Println(n)
	}

	// Output:
	// 1
	// 3
	// 5
}

func TestFilter(t *testing.T) {
	testCases := []struct {
		name      string
		input     []int
		predicate func(int) bool
		expected  []int
	}{
		{
			"empty_sequence",
			[]int{},
			func(n int) bool { return n%2 == 1 },
			[]int{},
		},
		{
			"filter_odd_numbers",
			[]int{1, 2, 3, 4, 5, 6},
			func(n int) bool { return n%2 == 1 },
			[]int{1, 3, 5},
		},
		{
			"filter_even_numbers",
			[]int{1, 2, 3, 4, 5, 6},
			func(n int) bool { return n%2 == 0 },
			[]int{2, 4, 6},
		},
		{
			"filter_greater_than_three",
			[]int{1, 2, 3, 4, 5},
			func(n int) bool { return n > 3 },
			[]int{4, 5},
		},
		{
			"no_matches",
			[]int{1, 2, 3},
			func(n int) bool { return n > 10 },
			[]int{},
		},
		{
			"all_matches",
			[]int{1, 2, 3},
			func(n int) bool { return n > 0 },
			[]int{1, 2, 3},
		},
		{
			"single_element_match",
			[]int{42},
			func(n int) bool { return n == 42 },
			[]int{42},
		},
		{
			"single_element_no_match",
			[]int{42},
			func(n int) bool { return n == 24 },
			[]int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := slices.Values(tc.input)
			actual := slices.Collect(seq.Filter(input, tc.predicate))

			if !slices.Equal(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func ExampleFilter2() {
	numbers := maps.All(map[string]int{"one": 1, "two": 2, "three": 3, "four": 4, "five": 5})

	oddAndLong := func(k string, v int) bool {
		return len(k) > 3 && v%2 == 1
	}

	oddAndLongNumbers := seq.Filter2(numbers, oddAndLong)

	printSorted(oddAndLongNumbers)

	// Output:
	// five: 5
	// three: 3
}

func TestFilter2(t *testing.T) {
	testCases := []struct {
		name      string
		input     iter.Seq2[string, int]
		predicate func(string, int) bool
		expected  map[string]int
	}{
		{
			"empty_sequence",
			maps.All(map[string]int{}),
			func(k string, v int) bool { return len(k) > 3 && v%2 == 1 },
			map[string]int{},
		},
		{
			"filter_odd_values_and_long_keys",
			maps.All(map[string]int{"one": 1, "two": 2, "three": 3, "four": 4, "five": 5}),
			func(k string, v int) bool { return len(k) > 3 && v%2 == 1 },
			map[string]int{"three": 3, "five": 5},
		},
		{
			"filter_even_values",
			maps.All(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}),
			func(k string, v int) bool { return v%2 == 0 },
			map[string]int{"b": 2, "d": 4},
		},
		{
			"filter_by_key_length",
			maps.All(map[string]int{"a": 1, "bb": 2, "ccc": 3, "dddd": 4}),
			func(k string, v int) bool { return len(k) >= 3 },
			map[string]int{"ccc": 3, "dddd": 4},
		},
		{
			"no_matches",
			maps.All(map[string]int{"a": 1, "b": 2}),
			func(k string, v int) bool { return v > 10 },
			map[string]int{},
		},
		{
			"all_matches",
			maps.All(map[string]int{"a": 1, "b": 2, "c": 3}),
			func(k string, v int) bool { return v > 0 },
			map[string]int{"a": 1, "b": 2, "c": 3},
		},
		{
			"single_pair_match",
			maps.All(map[string]int{"key": 42}),
			func(k string, v int) bool { return v == 42 },
			map[string]int{"key": 42},
		},
		{
			"single_pair_no_match",
			maps.All(map[string]int{"key": 42}),
			func(k string, v int) bool { return v == 24 },
			map[string]int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := maps.Collect(seq.Filter2(tc.input, tc.predicate))

			if !maps.Equal(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}
