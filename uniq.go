package seq

import "iter"

// Uniq ensures only unique values are returned from a sequence.
//
// Items are returned in the order they first appear.
func Uniq[V comparable](seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		seen := make(map[V]struct{})

		seq(func(v V) bool {
			if _, ok := seen[v]; ok {
				return true // already seen, skip
			}

			seen[v] = struct{}{}

			return yield(v)
		})
	}
}

// Uniq2 ensures only unique keys are returned from a sequence.
//
// Note: You can achieve a similar effect using [maps.Collect],
// but that requires loading the entire sequence into memory.
// Additionally, while [maps.Collect] returns the last seen value for each key,
// Uniq2 returns the first seen value for each key.
func Uniq2[K comparable, V any](seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		seen := make(map[K]struct{})

		seq(func(k K, v V) bool {
			if _, ok := seen[k]; ok {
				return true // already seen key, skip
			}

			seen[k] = struct{}{}

			return yield(k, v)
		})
	}
}
