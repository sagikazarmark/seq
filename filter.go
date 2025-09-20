package seq

import "iter"

// Filter returns a sequence that only yields values matching the predicate.
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

// Filter2 returns a sequence that only yields pairs matching the predicate.
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
