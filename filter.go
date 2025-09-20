package seq

import "iter"

// Filter creates an iterator using a predicate to determine if a value should be yielded.
//
// The returned iterator will yield only the values for which the predicate is true.
func Filter[V any](seq iter.Seq[V], predicate func(V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if !predicate(v) {
				continue
			}

			if !yield(v) {
				return
			}
		}
	}
}

// Filter2 creates an iterator using a predicate to determine if a pair should be yielded.
//
// The returned iterator will yield only the pairs for which the predicate is true.
func Filter2[K any, V any](seq iter.Seq2[K, V], predicate func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			if !predicate(k, v) {
				continue
			}

			if !yield(k, v) {
				return
			}
		}
	}
}
