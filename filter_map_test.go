package seq_test

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"testing"

	"github.com/sagikazarmark/seq"
)

func ExampleFilterMap() {
	numbers := slices.Values([]int{1, 2, 3, 4, 5, 6})

	// Double only odd numbers
	doubleIfOdd := func(n int) (int, bool) {
		return n * 2, n%2 == 1
	}

	doubled := seq.FilterMap(numbers, doubleIfOdd)

	for n := range doubled {
		fmt.Println(n)
	}

	// Output:
	// 2
	// 6
	// 10
}

func TestFilterMap(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		fn       func(int) (int, bool)
		expected []int
	}{
		{
			"empty_sequence",
			[]int{},
			func(n int) (int, bool) { return n * 2, n%2 == 1 },
			[]int{},
		},
		{
			"double_odd_numbers",
			[]int{1, 2, 3, 4, 5, 6},
			func(n int) (int, bool) { return n * 2, n%2 == 1 },
			[]int{2, 6, 10},
		},
		{
			"square_positive_numbers",
			[]int{-2, -1, 0, 1, 2, 3},
			func(n int) (int, bool) { return n * n, n > 0 },
			[]int{1, 4, 9},
		},
		{
			"no_matches",
			[]int{1, 2, 3},
			func(n int) (int, bool) { return n, n > 10 },
			[]int{},
		},
		{
			"all_matches",
			[]int{1, 2, 3},
			func(n int) (int, bool) { return n + 10, true },
			[]int{11, 12, 13},
		},
		{
			"single_element_match",
			[]int{5},
			func(n int) (int, bool) { return n * 3, n == 5 },
			[]int{15},
		},
		{
			"single_element_no_match",
			[]int{5},
			func(n int) (int, bool) { return n * 3, n == 10 },
			[]int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := slices.Values(tc.input)
			actual := slices.Collect(seq.FilterMap(input, tc.fn))

			if !slices.Equal(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func ExampleFilterMap2() {
	cart := maps.All(map[string]int{"apple": 5, "banana": 6, "dates": 4})
	prices := map[string]float64{"apple": 0.5, "banana": 0.4, "cherry": 0.2}

	// Filter out items with no price tag
	subtotal := func(k string, v int) (string, bool) {
		price, ok := prices[k]

		return fmt.Sprintf("$%.2f", float64(v)*price), ok
	}

	subtotals := seq.FilterMap2(cart, subtotal)

	printSorted(subtotals)

	// Output:
	// apple: $2.50
	// banana: $2.40
}

func TestFilterMap2(t *testing.T) {
	testCases := []struct {
		name     string
		input    iter.Seq2[string, int]
		fn       func(string, int) (int, bool)
		expected map[string]int
	}{
		{
			"empty_sequence",
			maps.All(map[string]int{}),
			func(k string, v int) (int, bool) { return v * 2, len(k) == v },
			map[string]int{},
		},
		{
			"filter_by_key_length_and_double",
			maps.All(map[string]int{"a": 1, "bb": 2, "ccc": 3, "dddd": 5}),
			func(k string, v int) (int, bool) { return v * 2, len(k) == v },
			map[string]int{"a": 2, "bb": 4, "ccc": 6},
		},
		{
			"filter_even_values_and_square",
			maps.All(map[string]int{"x": 2, "y": 3, "z": 4}),
			func(k string, v int) (int, bool) { return v * v, v%2 == 0 },
			map[string]int{"x": 4, "z": 16},
		},
		{
			"no_matches",
			maps.All(map[string]int{"a": 1, "b": 2}),
			func(k string, v int) (int, bool) { return v, v > 10 },
			map[string]int{},
		},
		{
			"all_matches",
			maps.All(map[string]int{"a": 1, "b": 2}),
			func(k string, v int) (int, bool) { return v + 10, true },
			map[string]int{"a": 11, "b": 12},
		},
		{
			"single_pair_match",
			maps.All(map[string]int{"key": 5}),
			func(k string, v int) (int, bool) { return v * 3, v == 5 },
			map[string]int{"key": 15},
		},
		{
			"single_pair_no_match",
			maps.All(map[string]int{"key": 5}),
			func(k string, v int) (int, bool) { return v * 3, v == 10 },
			map[string]int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := maps.Collect(seq.FilterMap2(tc.input, tc.fn))

			if !maps.Equal(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}
