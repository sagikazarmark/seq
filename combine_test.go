package seq_test

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/sagikazarmark/seq"
)

func ExampleCombine() {
	users1 := slices.Values([]string{"alice", "bob"})
	users2 := slices.Values([]string{"charlie", "dave"})

	users := seq.Combine(users1, users2)

	for user := range users {
		fmt.Println(user)
	}

	// Output:
	// alice
	// bob
	// charlie
	// dave
}

func TestCombine(t *testing.T) {
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
			actual := slices.Collect(seq.Combine(tc.seqs...))

			if !slices.Equal(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func ExampleCombine2() {
	roles1 := maps.All(map[string]string{"alice": "admin", "bob": "admin"})
	roles2 := maps.All(map[string]string{"bob": "user", "charlie": "manager"})

	roles := seq.Combine2(roles1, roles2)

	// Maps are unordered, so we need to hack here for this example
	var output []string

	for user, role := range roles {
		output = append(output, fmt.Sprintf("%s: %s\n", user, role))
	}

	sort.Strings(output)
	fmt.Print(strings.Join(output, ""))

	// Output:
	// alice: admin
	// bob: admin
	// bob: user
	// charlie: manager
}

func TestCombine2(t *testing.T) {
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
			actual := maps.Collect(seq.Combine2(tc.seqs...))

			if !maps.Equal(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}
