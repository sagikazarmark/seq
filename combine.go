package seq

import "iter"

// Combine merges multiple iterator sequences into a single sequence by
// concatenating them in order.
// Values from the first sequence are yielded first,
// followed by values from the second sequence, and so on.
func Combine[V any](seqs ...iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, seq := range seqs {
			for v := range seq {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Combine2 merges multiple iterator sequences into a single sequence by
// concatenating them in order.
// Pairs from the first sequence are yielded first,
// followed by pairs from the second sequence, and so on.
//
// If multiple sequences contain the same key, the later sequences will overwrite
// the earlier ones when collected into a map. However, when iterating directly,
// all key-value pairs are yielded in order, including duplicates.
func Combine2[K comparable, V any](seqs ...iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, seq := range seqs {
			for k, v := range seq {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}
