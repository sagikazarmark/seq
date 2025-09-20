package seq

import (
	"iter"
	"slices"
)

// ValuesErr returns an iterator that yields the slice elements in order, or an error
// if one was provided.
//
// If err is not nil, ValuesErr returns (nil, err). If err is nil, it returns
// an iterator over the slice values.
//
// This is useful for chaining operations without intermediate variables:
//
//	items, err := seq.ValuesErr(fetchItems())
//	if err != nil {
//		return err
//	}
//
//	for item := range items {
//		// process item
//	}
//
// Instead of:
//
//	itemSlice, err := fetchItems()
//	if err != nil {
//		return err
//	}
//
//	items := slices.Values(itemSlice)
//	for item := range items {
//		// process item
//	}
func ValuesErr[Slice ~[]E, E any](s Slice, err error) (iter.Seq[E], error) {
	if err != nil {
		return nil, err
	}

	return slices.Values(s), nil
}
