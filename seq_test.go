package seq_test

import (
	"fmt"
	"iter"
	"sort"
	"strings"
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
