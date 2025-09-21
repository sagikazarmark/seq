package seq_test

import (
	"errors"
	"fmt"
	"iter"
	"maps"
	"slices"
	"sort"
	"strings"

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

func ExampleSkip() {
	numbers := slices.Values([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	skip3 := seq.Skip(numbers, 3)

	for n := range skip3 {
		fmt.Println(n)
	}

	// Output:
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
	// 10
}

func ExampleSkip2() {
	data := maps.All(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5})

	skip2 := seq.Skip2(data, 2)

	count := 0
	for k, v := range skip2 {
		fmt.Printf("%s: %d\n", k, v)
		count++
	}
	fmt.Printf("Total pairs after skip: %d\n", count)

	// Output will vary due to map iteration order, but will show 3 pairs (5 - 2 = 3)
}

func ExampleTake() {
	numbers := slices.Values([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	first5 := seq.Take(numbers, 5)

	for n := range first5 {
		fmt.Println(n)
	}

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

func ExampleTake2() {
	data := maps.All(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5})

	first3 := seq.Take2(data, 3)

	count := 0
	for k, v := range first3 {
		fmt.Printf("%s: %d\n", k, v)
		count++
	}
	fmt.Printf("Total pairs: %d\n", count)

	// Output (map iteration order is not guaranteed, but we'll get 3 pairs):
	// Total pairs: 3
}

func ExampleTakeWhile() {
	numbers := slices.Values([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	lessThan5 := seq.TakeWhile(numbers, func(n int) bool { return n < 5 })

	for n := range lessThan5 {
		fmt.Println(n)
	}

	// Output:
	// 1
	// 2
	// 3
	// 4
}

func ExampleTakeWhile2() {
	data := maps.All(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5})

	// Note: Map iteration order is not guaranteed, so this example may vary
	smallValues := seq.TakeWhile2(data, func(k string, v int) bool { return v < 3 })

	count := 0
	for k, v := range smallValues {
		fmt.Printf("%s: %d\n", k, v)
		count++
	}
	fmt.Printf("Total pairs taken: %d\n", count)

	// Output will vary due to map iteration order, but will show pairs with values < 3
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
