package seq_test

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"testing"

	"github.com/sagikazarmark/seq"
)

func ExampleMap() {
	numbers := slices.Values([]int{1, 2, 3, 4, 5})

	double := func(n int) int {
		return n * 2
	}

	doubled := seq.Map(numbers, double)

	for n := range doubled {
		fmt.Println(n)
	}

	// Output:
	// 2
	// 4
	// 6
	// 8
	// 10
}

func TestMap(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		fn       func(int) int
		expected []int
	}{
		{
			"empty_sequence",
			[]int{},
			func(n int) int { return n * 2 },
			[]int{},
		},
		{
			"double_numbers",
			[]int{1, 2, 3, 4, 5},
			func(n int) int { return n * 2 },
			[]int{2, 4, 6, 8, 10},
		},
		{
			"square_numbers",
			[]int{1, 2, 3, 4},
			func(n int) int { return n * n },
			[]int{1, 4, 9, 16},
		},
		{
			"add_ten",
			[]int{5, 10, 15},
			func(n int) int { return n + 10 },
			[]int{15, 20, 25},
		},
		{
			"negate_numbers",
			[]int{1, -2, 3, -4},
			func(n int) int { return -n },
			[]int{-1, 2, -3, 4},
		},
		{
			"single_element",
			[]int{42},
			func(n int) int { return n * 3 },
			[]int{126},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := slices.Values(tc.input)
			actual := slices.Collect(seq.Map(input, tc.fn))

			if !slices.Equal(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func ExampleMap2() {
	cart := maps.All(map[string]int{"apple": 5, "banana": 6, "cherry": 4})
	prices := map[string]float64{"apple": 0.5, "banana": 0.4, "cherry": 0.2}

	subtotal := func(k string, v int) string {
		return fmt.Sprintf("$%.2f", float64(v)*prices[k])
	}

	subtotals := seq.Map2(cart, subtotal)

	printSorted(subtotals)

	// Output:
	// apple: $2.50
	// banana: $2.40
	// cherry: $0.80
}

func TestMap2(t *testing.T) {
	testCases := []struct {
		name     string
		input    iter.Seq2[string, int]
		fn       func(string, int) int
		expected map[string]int
	}{
		{
			"empty_sequence",
			maps.All(map[string]int{}),
			func(k string, v int) int { return v * 2 },
			map[string]int{},
		},
		{
			"double_values",
			maps.All(map[string]int{"a": 1, "b": 2, "c": 3}),
			func(k string, v int) int { return v * 2 },
			map[string]int{"a": 2, "b": 4, "c": 6},
		},
		{
			"add_key_length_to_value",
			maps.All(map[string]int{"a": 1, "bb": 2, "ccc": 3}),
			func(k string, v int) int { return v + len(k) },
			map[string]int{"a": 2, "bb": 4, "ccc": 6},
		},
		{
			"multiply_by_key_length",
			maps.All(map[string]int{"x": 5, "yy": 3}),
			func(k string, v int) int { return v * len(k) },
			map[string]int{"x": 5, "yy": 6},
		},
		{
			"single_pair",
			maps.All(map[string]int{"key": 10}),
			func(k string, v int) int { return v + 5 },
			map[string]int{"key": 15},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := maps.Collect(seq.Map2(tc.input, tc.fn))

			if !maps.Equal(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}
