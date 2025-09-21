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
	users1 := maps.All(map[string]string{"alice": "admin", "bob": "admin"})
	users2 := maps.All(map[string]string{"bob": "user", "charlie": "manager"})

	users := seq.Chain2(users1, users2)

	printSorted(users)

	// Output:
	// alice: admin
	// bob: admin
	// bob: user
	// charlie: manager
}

func ExampleFilter() {
	numbers := slices.Values([]int{1, 2, 3, 4, 5})

	oddFilter := func(n int) bool {
		return n%2 == 1
	}

	oddNumbers := seq.Filter(numbers, oddFilter)

	for n := range oddNumbers {
		fmt.Println(n)
	}

	// Output:
	// 1
	// 3
	// 5
}

func ExampleFilter2() {
	users := maps.All(map[string]string{"alice": "admin", "bob": "user", "charlie": "manager", "dave": "manager"})

	managerFilter := func(k string, v string) bool {
		return len(k) > 3 && v == "manager"
	}

	managers := seq.Filter2(users, managerFilter)

	printSorted(managers)

	// Output:
	// charlie: manager
	// dave: manager
}

func ExampleFilterMap() {
	numbers := slices.Values([]int{1, 2, 3, 4, 5, 6})

	// Filter odd numbers and double them
	doubleIfOdd := func(n int) (int, bool) {
		return n * 2, n%2 == 1
	}

	doubledOdd := seq.FilterMap(numbers, doubleIfOdd)

	for n := range doubledOdd {
		fmt.Println(n)
	}

	// Output:
	// 2
	// 6
	// 10
}

func ExampleFilterMap2() {
	salaries := maps.All(map[string]float64{"alice": 100000, "bob": 80000, "charlie": 120000, "dave": 50000})

	// Give everyone with a salary under 100,000 credits a raise of 5%
	// and produce a list of names and new salaries
	raise := func(k string, v float64) (float64, bool) {
		if v < 100000 {
			return v * 1.05, true
		}

		return v, false
	}

	promotions := seq.FilterMap2(salaries, raise)

	printSorted(promotions)

	// Output:
	// bob: 84000
	// dave: 52500
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
	users := maps.All(map[string]string{"alice": "admin", "bob": "user", "charlie": "manager", "dave": "manager"})
	payGrades := map[string]int{"admin": 1000, "user": 500, "manager": 1500}

	// Calculate salaries for each user based on their role
	calculateSalary := func(k string, v string) int {
		return payGrades[v]
	}

	salaries := seq.Map2(users, calculateSalary)

	printSorted(salaries)

	// Output:
	// alice: 1000
	// bob: 500
	// charlie: 1500
	// dave: 1500
}

func ExampleRepeat() {
	meaningOfLife := seq.Repeat(42)

	var i int

	for n := range meaningOfLife {
		fmt.Println(n)
		i++
		if i == 5 {
			break
		}
	}

	// Output:
	// 42
	// 42
	// 42
	// 42
	// 42
}

func ExampleSkip() {
	fruits := slices.Values([]string{"apple", "banana", "cherry", "grape", "mango"})

	skip3 := seq.Skip(fruits, 3)

	for n := range skip3 {
		fmt.Println(n)
	}

	// Output:
	// grape
	// mango
}

func ExampleSkip2() {
	users := seq.Sorted2(map[string]string{"alice": "admin", "bob": "user", "charlie": "manager", "dave": "manager"})

	skip2 := seq.Skip2(users, 2)

	printSorted(skip2)

	// Output:
	// charlie: manager
	// dave: manager
}

func ExampleSkipWhile() {
	fruits := slices.Values([]string{"apple", "apricot", "acerola", "banana", "cherry", "grape", "mango"})

	noAs := seq.SkipWhile(fruits, func(v string) bool { return v[0] == 'a' })

	for n := range noAs {
		fmt.Println(n)
	}

	// Output:
	// banana
	// cherry
	// grape
	// mango
}

func ExampleSkipWhile2() {
	users := seq.Sorted2(map[string]string{"alice": "admin", "bob": "admin", "charlie": "manager", "dave": "user"})

	skipAdmins := func(k string, v string) bool {
		return v == "admin"
	}

	afterAdmins := seq.SkipWhile2(users, skipAdmins)

	printSorted(afterAdmins)

	// Output:
	// charlie: manager
	// dave: user
}

func ExampleSorted2() {
	users := seq.Sorted2(map[string]string{"charlie": "manager", "bob": "user", "alice": "admin", "dave": "manager"})

	for user, role := range users {
		fmt.Printf("%s: %s\n", user, role)
	}

	// Output:
	// alice: admin
	// bob: user
	// charlie: manager
	// dave: manager
}

func ExampleTake() {
	fruits := slices.Values([]string{"apple", "banana", "cherry", "grape", "mango"})

	first3 := seq.Take(fruits, 3)

	for n := range first3 {
		fmt.Println(n)
	}

	// Output:
	// apple
	// banana
	// cherry
}

func ExampleTake2() {
	users := seq.Sorted2(map[string]string{"alice": "admin", "bob": "user", "charlie": "manager", "dave": "manager"})

	first3 := seq.Take2(users, 3)

	printSorted(first3)

	// Output:
	// alice: admin
	// bob: user
	// charlie: manager
}

func ExampleTakeWhile() {
	fruits := slices.Values([]string{"apple", "apricot", "acerola", "banana", "cherry", "grape", "mango"})

	startsWithA := seq.TakeWhile(fruits, func(v string) bool { return v[0] == 'a' })

	for n := range startsWithA {
		fmt.Println(n)
	}

	// Output:
	// apple
	// apricot
	// acerola
}

func ExampleTakeWhile2() {
	users := seq.Sorted2(map[string]string{"alice": "admin", "bob": "admin", "charlie": "manager", "dave": "manager"})

	firstAdmins := seq.TakeWhile2(users, func(k string, v string) bool { return v == "admin" })

	printSorted(firstAdmins)

	// Output:
	// alice: admin
	// bob: admin
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

// Maps are unordered, so we need to hack for examples
func printSorted[K any, V any](s iter.Seq2[K, V]) {
	var output []string

	for key, value := range s {
		output = append(output, fmt.Sprintf("%v: %v\n", key, value))
	}

	sort.Strings(output)
	fmt.Print(strings.Join(output, ""))
}
