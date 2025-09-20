package seq

import "iter"

// Flatten returns a sequence that yields all elements of each slice
// produced by the outer sequence.
func Flatten[V any](seq iter.Seq[[]V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		seq(func(slice []V) bool {
			for _, v := range slice {
				if !yield(v) {
					return false
				}
			}

			return true
		})
	}
}
