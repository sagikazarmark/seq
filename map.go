package seq

import "iter"

// Map creates an iterator that transforms values using a function.
//
// The returned iterator will yield the transformed values.
func Map[V any, U any](seq iter.Seq[V], fn func(V) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for v := range seq {
			if !yield(fn(v)) {
				return
			}
		}
	}
}

// Map2 creates an iterator that transforms values of a pair using a function.
//
// The returned iterator will yield the transformed values.
func Map2[K any, V any, U any](seq iter.Seq2[K, V], fn func(K, V) U) iter.Seq2[K, U] {
	return func(yield func(K, U) bool) {
		for k, v := range seq {
			if !yield(k, fn(k, v)) {
				return
			}
		}
	}
}
