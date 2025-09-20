package seq_test

import (
	"errors"
	"fmt"
	"iter"
	"maps"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/sagikazarmark/seq"
)

// Maps are unordered, so we need to hack for examples
func printSorted[K any, V any](s iter.Seq2[K, V]) {
	var output []string

	for key, value := range s {
		output = append(output, fmt.Sprintf("%v: %v\n", key, value))
	}

	sort.Strings(output)
	fmt.Print(strings.Join(output, ""))
}

func ExampleChain() {
	users1 := slices.Values([]string{"alice", "bob"})
	users2 := slices.Values([]string{"charlie", "dave"})

	users := seq.Chain(users1, users2)

	for user := range users {
		fmt.Println(user)
	}

	// Output:
	// alice
	// bob
	// charlie
	// dave
}

func TestChain(t *testing.T) {
	testCases := []struct {
		name     string
		seqs     []iter.Seq[int]
		expected []int
	}{
		{"no_sequences", []iter.Seq[int]{}, []int{}},
		{"empty_sequences", []iter.Seq[int]{slices.Values([]int{}), slices.Values([]int{})}, []int{}},
		{"single_sequence", []iter.Seq[int]{slices.Values([]int{1, 2, 3})}, []int{1, 2, 3}},
		{"multiple_sequences", []iter.Seq[int]{slices.Values([]int{1, 2}), slices.Values([]int{3, 4}), slices.Values([]int{5})}, []int{1, 2, 3, 4, 5}},
		{"empty_sequence_mixed_with_non_empty", []iter.Seq[int]{slices.Values([]int{1, 2}), slices.Values([]int{}), slices.Values([]int{3})}, []int{1, 2, 3}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := slices.Collect(seq.Chain(tc.seqs...))

			if !slices.Equal(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func ExampleChain2() {
	roles1 := maps.All(map[string]string{"alice": "admin", "bob": "admin"})
	roles2 := maps.All(map[string]string{"bob": "user", "charlie": "manager"})

	roles := seq.Chain2(roles1, roles2)

	printSorted(roles)

	// Output:
	// alice: admin
	// bob: admin
	// bob: user
	// charlie: manager
}

func TestChain2(t *testing.T) {
	testCases := []struct {
		name     string
		seqs     []iter.Seq2[string, int]
		expected map[string]int
	}{
		{"empty_sequences", []iter.Seq2[string, int]{}, map[string]int{}},
		{"single_sequence", []iter.Seq2[string, int]{maps.All(map[string]int{"a": 1, "b": 2})}, map[string]int{"a": 1, "b": 2}},
		{"multiple_sequences", []iter.Seq2[string, int]{maps.All(map[string]int{"a": 1, "b": 2}), maps.All(map[string]int{"c": 3, "d": 4})}, map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}},
		{"empty_sequence_mixed_with_non_empty", []iter.Seq2[string, int]{maps.All(map[string]int{"a": 1}), maps.All(map[string]int{}), maps.All(map[string]int{"b": 2})}, map[string]int{"a": 1, "b": 2}},
		{"duplicate_keys_from_different_sequences", []iter.Seq2[string, int]{maps.All(map[string]int{"a": 1, "b": 2}), maps.All(map[string]int{"b": 3, "c": 4})}, map[string]int{"a": 1, "b": 3, "c": 4}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := maps.Collect(seq.Chain2(tc.seqs...))

			if !maps.Equal(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

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

func ExampleValuesErr_ok() {
	fn := func() ([]string, error) {
		return []string{"apple", "banana", "cherry"}, nil
	}

	fruits, err := seq.ValuesErr(fn())
	if err != nil {
		panic(err)
	}

	for fruit := range fruits {
		fmt.Println(fruit)
	}

	// Output:
	// apple
	// banana
	// cherry
}

func ExampleValuesErr_error() {
	fn := func() ([]struct{}, error) {
		return nil, errors.New("something went wrong")
	}

	_, err := seq.ValuesErr(fn())
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// something went wrong
}
