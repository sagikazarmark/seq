package seq_test

import (
	"fmt"
	"maps"
	"slices"

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

	printSorted(oddAndLongNumbers)

	// Output:
	// five: 5
	// three: 3
}
