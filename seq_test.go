package seq_test

import (
	"iter"
	"maps"
	"slices"
	"testing"

	"github.com/sagikazarmark/seq"
)

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

func TestRepeat(t *testing.T) {
	testCases := []struct {
		name     string
		value    int
		takeN    uint
		expected []int
	}{
		{
			"repeat_integer",
			42,
			5,
			[]int{42, 42, 42, 42, 42},
		},
		{
			"repeat_zero_times",
			10,
			0,
			[]int{},
		},
		{
			"repeat_once",
			7,
			1,
			[]int{7},
		},
		{
			"repeat_negative_number",
			-5,
			3,
			[]int{-5, -5, -5},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repeated := seq.Repeat(tc.value)
			actual := slices.Collect(seq.Take(repeated, tc.takeN))

			if !slices.Equal(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestSkip(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		n        uint
		expected []int
	}{
		{
			"empty_sequence",
			[]int{},
			3,
			[]int{},
		},
		{
			"skip_zero",
			[]int{1, 2, 3, 4, 5},
			0,
			[]int{1, 2, 3, 4, 5},
		},
		{
			"skip_less_than_available",
			[]int{1, 2, 3, 4, 5},
			3,
			[]int{4, 5},
		},
		{
			"skip_exactly_available",
			[]int{1, 2, 3, 4, 5},
			5,
			[]int{},
		},
		{
			"skip_more_than_available",
			[]int{1, 2, 3},
			10,
			[]int{},
		},
		{
			"skip_one_from_single_element",
			[]int{42},
			1,
			[]int{},
		},
		{
			"skip_zero_from_single_element",
			[]int{42},
			0,
			[]int{42},
		},
		{
			"skip_more_than_single_element",
			[]int{42},
			5,
			[]int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := slices.Values(tc.input)
			actual := slices.Collect(seq.Skip(input, tc.n))

			if !slices.Equal(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestSkip2(t *testing.T) {
	testCases := []struct {
		name        string
		input       map[string]int
		n           uint
		expectedLen int
	}{
		{
			"empty_sequence",
			map[string]int{},
			3,
			0,
		},
		{
			"skip_zero",
			map[string]int{"a": 1, "b": 2, "c": 3},
			0,
			3,
		},
		{
			"skip_less_than_available",
			map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5},
			2,
			3,
		},
		{
			"skip_exactly_available",
			map[string]int{"a": 1, "b": 2, "c": 3},
			3,
			0,
		},
		{
			"skip_more_than_available",
			map[string]int{"a": 1, "b": 2},
			10,
			0,
		},
		{
			"skip_one_from_single_pair",
			map[string]int{"key": 42},
			1,
			0,
		},
		{
			"skip_zero_from_single_pair",
			map[string]int{"key": 42},
			0,
			1,
		},
		{
			"skip_more_than_single_pair",
			map[string]int{"key": 42},
			5,
			0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := maps.All(tc.input)
			actual := maps.Collect(seq.Skip2(input, tc.n))

			if len(actual) != tc.expectedLen {
				t.Errorf("expected length %d, got %d", tc.expectedLen, len(actual))
			}

			// Verify that all returned pairs exist in the original input
			for k, v := range actual {
				if originalV, exists := tc.input[k]; !exists || originalV != v {
					t.Errorf("unexpected pair %s: %d not in original input", k, v)
				}
			}
		})
	}
}

func TestSkipWhile(t *testing.T) {
	testCases := []struct {
		name      string
		input     []int
		predicate func(int) bool
		expected  []int
	}{
		{
			"empty_sequence",
			[]int{},
			func(n int) bool { return n < 5 },
			[]int{},
		},
		{
			"nil_predicate",
			[]int{1, 2, 3, 4, 5},
			nil,
			[]int{1, 2, 3, 4, 5},
		},
		{
			"skip_while_less_than_5",
			[]int{1, 2, 3, 4, 5, 6, 7, 3, 4},
			func(n int) bool { return n < 5 },
			[]int{5, 6, 7, 3, 4},
		},
		{
			"skip_while_even",
			[]int{2, 4, 6, 8, 9, 10, 12},
			func(n int) bool { return n%2 == 0 },
			[]int{9, 10, 12},
		},
		{
			"predicate_never_true",
			[]int{1, 2, 3, 4, 5},
			func(n int) bool { return n > 10 },
			[]int{1, 2, 3, 4, 5},
		},
		{
			"predicate_always_true",
			[]int{1, 2, 3, 4, 5},
			func(n int) bool { return n > 0 },
			[]int{},
		},
		{
			"single_element_true",
			[]int{42},
			func(n int) bool { return n == 42 },
			[]int{},
		},
		{
			"single_element_false",
			[]int{42},
			func(n int) bool { return n != 42 },
			[]int{42},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := slices.Values(tc.input)
			actual := slices.Collect(seq.SkipWhile(input, tc.predicate))

			if !slices.Equal(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestSkipWhile2(t *testing.T) {
	testCases := []struct {
		name              string
		input             map[string]int
		predicate         func(string, int) bool
		expectedCondition func(map[string]int) bool
		description       string
	}{
		{
			"empty_sequence",
			map[string]int{},
			func(k string, v int) bool { return v < 5 },
			func(result map[string]int) bool { return len(result) == 0 },
			"should return empty map",
		},
		{
			"nil_predicate",
			map[string]int{"a": 1, "b": 2},
			nil,
			func(result map[string]int) bool { return len(result) == 2 },
			"should return all pairs when predicate is nil",
		},
		{
			"predicate_never_true",
			map[string]int{"a": 1, "b": 2, "c": 3},
			func(k string, v int) bool { return v > 10 },
			func(result map[string]int) bool { return len(result) == 3 },
			"should return all pairs when predicate never matches",
		},
		{
			"single_element_true",
			map[string]int{"key": 42},
			func(k string, v int) bool { return v == 42 },
			func(result map[string]int) bool { return len(result) == 0 },
			"should return empty when single element matches",
		},
		{
			"single_element_false",
			map[string]int{"key": 42},
			func(k string, v int) bool { return v != 42 },
			func(result map[string]int) bool {
				return len(result) == 1 && result["key"] == 42
			},
			"should return single element when it doesn't match",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := maps.All(tc.input)
			actual := maps.Collect(seq.SkipWhile2(input, tc.predicate))

			if !tc.expectedCondition(actual) {
				t.Errorf("%s, got %v", tc.description, actual)
			}

			// Verify that all returned pairs exist in the original input
			for k, v := range actual {
				if originalV, exists := tc.input[k]; !exists || originalV != v {
					t.Errorf("unexpected pair %s: %d not in original input", k, v)
				}
			}
		})
	}
}

func TestSorted2(t *testing.T) {
	testCases := []struct {
		name        string
		input       map[string]int
		expectedSeq []struct {
			key   string
			value int
		}
	}{
		{
			"empty_map",
			map[string]int{},
			[]struct {
				key   string
				value int
			}{},
		},
		{
			"single_pair",
			map[string]int{"key": 42},
			[]struct {
				key   string
				value int
			}{{"key", 42}},
		},
		{
			"multiple_pairs_sorted",
			map[string]int{"charlie": 3, "alice": 1, "bob": 2},
			[]struct {
				key   string
				value int
			}{
				{"alice", 1},
				{"bob", 2},
				{"charlie", 3},
			},
		},
		{
			"numeric_keys",
			map[string]int{"10": 10, "2": 2, "1": 1},
			[]struct {
				key   string
				value int
			}{
				{"1", 1},
				{"10", 10},
				{"2", 2},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sortedSeq := seq.Sorted2(tc.input)

			var actual []struct {
				key   string
				value int
			}

			for k, v := range sortedSeq {
				actual = append(actual, struct {
					key   string
					value int
				}{k, v})
			}

			if len(actual) != len(tc.expectedSeq) {
				t.Errorf("expected length %d, got %d", len(tc.expectedSeq), len(actual))
				return
			}

			for i, expected := range tc.expectedSeq {
				if actual[i].key != expected.key || actual[i].value != expected.value {
					t.Errorf("at index %d: expected %+v, got %+v", i, expected, actual[i])
				}
			}
		})
	}
}

func TestTake(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		n        uint
		expected []int
	}{
		{
			"empty_sequence",
			[]int{},
			5,
			[]int{},
		},
		{
			"take_zero",
			[]int{1, 2, 3, 4, 5},
			0,
			[]int{},
		},
		{
			"take_less_than_available",
			[]int{1, 2, 3, 4, 5},
			3,
			[]int{1, 2, 3},
		},
		{
			"take_exactly_available",
			[]int{1, 2, 3, 4, 5},
			5,
			[]int{1, 2, 3, 4, 5},
		},
		{
			"take_more_than_available",
			[]int{1, 2, 3},
			10,
			[]int{1, 2, 3},
		},
		{
			"take_one_from_single_element",
			[]int{42},
			1,
			[]int{42},
		},
		{
			"take_zero_from_single_element",
			[]int{42},
			0,
			[]int{},
		},
		{
			"take_more_than_single_element",
			[]int{42},
			5,
			[]int{42},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := slices.Values(tc.input)
			actual := slices.Collect(seq.Take(input, tc.n))

			if !slices.Equal(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestTake2(t *testing.T) {
	testCases := []struct {
		name        string
		input       map[string]int
		n           uint
		expectedLen int
	}{
		{
			"empty_sequence",
			map[string]int{},
			5,
			0,
		},
		{
			"take_zero",
			map[string]int{"a": 1, "b": 2, "c": 3},
			0,
			0,
		},
		{
			"take_less_than_available",
			map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5},
			3,
			3,
		},
		{
			"take_exactly_available",
			map[string]int{"a": 1, "b": 2, "c": 3},
			3,
			3,
		},
		{
			"take_more_than_available",
			map[string]int{"a": 1, "b": 2},
			10,
			2,
		},
		{
			"take_one_from_single_pair",
			map[string]int{"key": 42},
			1,
			1,
		},
		{
			"take_zero_from_single_pair",
			map[string]int{"key": 42},
			0,
			0,
		},
		{
			"take_more_than_single_pair",
			map[string]int{"key": 42},
			5,
			1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := maps.All(tc.input)
			actual := maps.Collect(seq.Take2(input, tc.n))

			if len(actual) != tc.expectedLen {
				t.Errorf("expected length %d, got %d", tc.expectedLen, len(actual))
			}

			// Verify that all returned pairs exist in the original input
			for k, v := range actual {
				if originalV, exists := tc.input[k]; !exists || originalV != v {
					t.Errorf("unexpected pair %s: %d not in original input", k, v)
				}
			}
		})
	}
}

func TestTakeWhile(t *testing.T) {
	testCases := []struct {
		name      string
		input     []int
		predicate func(int) bool
		expected  []int
	}{
		{
			"empty_sequence",
			[]int{},
			func(n int) bool { return n < 5 },
			[]int{},
		},
		{
			"nil_predicate",
			[]int{1, 2, 3, 4, 5},
			nil,
			[]int{},
		},
		{
			"take_while_less_than_5",
			[]int{1, 2, 3, 4, 5, 6, 7},
			func(n int) bool { return n < 5 },
			[]int{1, 2, 3, 4},
		},
		{
			"take_while_even",
			[]int{2, 4, 6, 8, 9, 10, 12},
			func(n int) bool { return n%2 == 0 },
			[]int{2, 4, 6, 8},
		},
		{
			"predicate_never_true",
			[]int{1, 2, 3, 4, 5},
			func(n int) bool { return n > 10 },
			[]int{},
		},
		{
			"predicate_always_true",
			[]int{1, 2, 3, 4, 5},
			func(n int) bool { return n > 0 },
			[]int{1, 2, 3, 4, 5},
		},
		{
			"single_element_true",
			[]int{42},
			func(n int) bool { return n == 42 },
			[]int{42},
		},
		{
			"single_element_false",
			[]int{42},
			func(n int) bool { return n != 42 },
			[]int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := slices.Values(tc.input)
			actual := slices.Collect(seq.TakeWhile(input, tc.predicate))

			if !slices.Equal(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestTakeWhile2(t *testing.T) {
	testCases := []struct {
		name              string
		input             map[string]int
		predicate         func(string, int) bool
		expectedCondition func(map[string]int) bool
		description       string
	}{
		{
			"empty_sequence",
			map[string]int{},
			func(k string, v int) bool { return v < 5 },
			func(result map[string]int) bool { return len(result) == 0 },
			"should return empty map",
		},
		{
			"nil_predicate",
			map[string]int{"a": 1, "b": 2},
			nil,
			func(result map[string]int) bool { return len(result) == 0 },
			"should return empty map when predicate is nil",
		},
		{
			"predicate_never_true",
			map[string]int{"a": 1, "b": 2, "c": 3},
			func(k string, v int) bool { return v > 10 },
			func(result map[string]int) bool { return len(result) == 0 },
			"should return empty map when predicate never matches",
		},
		{
			"single_element_true",
			map[string]int{"key": 42},
			func(k string, v int) bool { return v == 42 },
			func(result map[string]int) bool {
				return len(result) == 1 && result["key"] == 42
			},
			"should return single element when it matches",
		},
		{
			"single_element_false",
			map[string]int{"key": 42},
			func(k string, v int) bool { return v != 42 },
			func(result map[string]int) bool { return len(result) == 0 },
			"should return empty when single element doesn't match",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := maps.All(tc.input)
			actual := maps.Collect(seq.TakeWhile2(input, tc.predicate))

			if !tc.expectedCondition(actual) {
				t.Errorf("%s, got %v", tc.description, actual)
			}

			// Verify that all returned pairs exist in the original input
			for k, v := range actual {
				if originalV, exists := tc.input[k]; !exists || originalV != v {
					t.Errorf("unexpected pair %s: %d not in original input", k, v)
				}
			}
		})
	}
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
