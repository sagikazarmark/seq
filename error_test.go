package seq_test

import (
	"errors"
	"fmt"

	"github.com/sagikazarmark/seq"
)

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
