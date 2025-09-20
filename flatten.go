package seq

import "iter"

// Flatten returns a sequence that yields all elements of each slice
// produced by the outer sequence. It takes a sequence of slices and
// "flattens" it into a sequence of individual elements.
//
// The elements are yielded in order: first all elements from the first slice,
// then all elements from the second slice, and so on.
//
// This is useful for processing nested data structures or combining
// results from multiple operations that return slices.
func Flatten[V any](seq iter.Seq[[]V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for slice := range seq {
			for _, v := range slice {
				if !yield(v) {
					return
				}
			}
		}
	}
}
