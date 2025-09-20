package seq_test

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"

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

func ExampleFilter2() {
	numbers := maps.All(map[string]int{"one": 1, "two": 2, "three": 3, "four": 4, "five": 5})

	oddAndLong := func(k string, v int) bool {
		return len(k) > 3 && v%2 == 1
	}

	oddAndLongNumbers := seq.Filter2(numbers, oddAndLong)

	// Maps are unordered, so we need to hack here for this example
	var output []string

	for s, i := range oddAndLongNumbers {
		output = append(output, fmt.Sprintf("%s: %d\n", s, i))
	}

	sort.Strings(output)
	fmt.Print(strings.Join(output, ""))

	// Output:
	// five: 5
	// three: 3
}
