package seq

import (
	"iter"
	"slices"
)

// ValuesErr returns an iterator that yields the slice elements in order or an error.
//
// It is useful to spare an intermediate variable when it's going to be used as an iterator.
func ValuesErr[Slice ~[]E, E any](s Slice, err error) (iter.Seq[E], error) {
	if err != nil {
		return nil, err
	}

	return slices.Values(s), nil
}
