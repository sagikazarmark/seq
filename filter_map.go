package seq

import "iter"

// FilterMap creates an iterator that both filters and maps.
//
// See [Filter] and [Map] for details.
func FilterMap[V any, U any](seq iter.Seq[V], fn func(V) (U, bool)) iter.Seq[U] {
	return func(yield func(U) bool) {
		for v := range seq {
			u, ok := fn(v)
			if !ok {
				continue
			}

			if !yield(u) {
				return
			}
		}
	}
}

// FilterMap2 creates an iterator that both filters and maps.
//
// See [Filter2] and [Map2] for details.
func FilterMap2[K any, V any, U any](seq iter.Seq2[K, V], fn func(K, V) (U, bool)) iter.Seq2[K, U] {
	return func(yield func(K, U) bool) {
		for k, v := range seq {
			u, ok := fn(k, v)
			if !ok {
				continue
			}

			if !yield(k, u) {
				return
			}
		}
	}
}
