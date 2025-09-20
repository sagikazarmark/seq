package seq

import "iter"

// Chain takes two (or more) iterators and creates a new iterator over them in sequence.
//
// The new iterator will iterate over values from each iterator in the order they are provided.
//
// In other words, it links two iterators together, in a chain.
func Chain[V any](seqs ...iter.Seq[V]) iter.Seq[V] {
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

// Chain2 takes two (or more) iterators and creates a new iterator over them in sequence.
//
// The new iterator will iterate over pairs from each iterator in the order they are provided.
//
// In other words, it links two iterators together, in a chain.
func Chain2[K any, V any](seqs ...iter.Seq2[K, V]) iter.Seq2[K, V] {
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
