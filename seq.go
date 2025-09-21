// Package seq provides utilities for working with Go 1.23+ iterator sequences.
package seq

import (
	"cmp"
	"iter"
	"maps"
	"slices"
)

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

// Flatten returns an iterator that yields all elements of each slice
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

// Repeat creates an iterator that yields the same value over and over.
//
// WARNING: This iterator will never terminate on its own.
func Repeat[V any](v V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for {
			if !yield(v) {
				return
			}
		}
	}
}

// Skip creates an iterator that skips values until n values are skipped or the end of the iterator is reached (whichever happens first).
func Skip[V comparable](seq iter.Seq[V], n uint) iter.Seq[V] {
	// Return early if n is zero
	if n == 0 {
		return seq
	}

	return func(yield func(V) bool) {
		var i uint

		for v := range seq {
			if i < n {
				i++

				continue
			}

			if !yield(v) {
				return
			}
		}
	}
}

// Skip2 creates an iterator that skips pairs until n pairs are skipped or the end of the iterator is reached (whichever happens first).
func Skip2[K comparable, V any](seq iter.Seq2[K, V], n uint) iter.Seq2[K, V] {
	// Return early if n is zero
	if n == 0 {
		return seq
	}

	return func(yield func(K, V) bool) {
		var i uint

		for k, v := range seq {
			if i < n {
				i++

				continue
			}

			if !yield(k, v) {
				return
			}
		}
	}
}

// SkipWhile creates an iterator that yields values based on a predicate.
//
// It will evaluate the predicate on each value, and skip values while it evaluates true.
// After false is returned, the rest of the values are yielded.
func SkipWhile[V comparable](seq iter.Seq[V], predicate func(V) bool) iter.Seq[V] {
	// Return early if predicate is nil
	if predicate == nil {
		return seq
	}

	return func(yield func(V) bool) {
		for v := range seq {
			if predicate(v) {
				continue
			}

			predicate = func(V) bool { return false }

			if !yield(v) {
				return
			}
		}
	}
}

// SkipWhile2 creates an iterator that yields pairs based on a predicate.
//
// It will evaluate the predicate on each pair, and skip pairs while it evaluates true.
// After false is returned, the rest of the pairs are yielded.
func SkipWhile2[K comparable, V any](seq iter.Seq2[K, V], predicate func(K, V) bool) iter.Seq2[K, V] {
	// Return early if predicate is nil
	if predicate == nil {
		return seq
	}

	return func(yield func(K, V) bool) {
		for k, v := range seq {
			if predicate(k, v) {
				continue
			}

			predicate = func(K, V) bool { return false }

			if !yield(k, v) {
				return
			}
		}
	}
}

// Sorted2 creates an iterator that yields pairs in a sorted order (by key).
//
// Useful when you need to iterate over a map in a sorted order.
func Sorted2[K cmp.Ordered, V any](m map[K]V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		keys := slices.Sorted(maps.Keys(m))

		for _, k := range keys {
			if !yield(k, m[k]) {
				return
			}
		}
	}
}

// Take creates an iterator that yields the first n values, or fewer if the underlying iterator ends sooner.
func Take[V comparable](seq iter.Seq[V], n uint) iter.Seq[V] {
	return func(yield func(V) bool) {
		// Return early if n is zero
		if n == 0 {
			return
		}

		var i uint

		for v := range seq {
			if !yield(v) {
				return
			}

			i++

			if i == n {
				return
			}
		}
	}
}

// Take2 creates an iterator that yields the first n pairs, or fewer if the underlying iterator ends sooner.
func Take2[K comparable, V any](seq iter.Seq2[K, V], n uint) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		// Return early if n is zero
		if n == 0 {
			return
		}

		var i uint

		for k, v := range seq {
			if !yield(k, v) {
				return
			}

			i++

			if i == n {
				return
			}
		}
	}
}

// TakeWhile creates an iterator that yields values based on a predicate.
//
// It will evaluate the predicate on each value, and yield values while it evaluates true.
// After false is returned, the rest of the values are ignored.
func TakeWhile[V comparable](seq iter.Seq[V], predicate func(V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		// Return early if predicate is nil
		if predicate == nil {
			return
		}

		for v := range seq {
			if !predicate(v) {
				return
			}

			if !yield(v) {
				return
			}
		}
	}
}

// TakeWhile2 creates an iterator that yields pairs based on a predicate.
//
// It will evaluate the predicate on each pair, and yield pairs while it evaluates true.
// After false is returned, the rest of the pairs are ignored.
func TakeWhile2[K comparable, V any](seq iter.Seq2[K, V], predicate func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		// Return early if predicate is nil
		if predicate == nil {
			return
		}

		for k, v := range seq {
			if !predicate(k, v) {
				return
			}

			if !yield(k, v) {
				return
			}
		}
	}
}

// Uniq ensures only unique values are returned from a sequence.
//
// Items are returned in the order they first appear.
func Uniq[V comparable](seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		seen := make(map[V]struct{})

		for v := range seq {
			if _, ok := seen[v]; ok {
				continue // already seen, skip
			}

			seen[v] = struct{}{}

			if !yield(v) {
				return
			}
		}
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

		for k, v := range seq {
			if _, ok := seen[k]; ok {
				continue // already seen key, skip
			}

			seen[k] = struct{}{}

			if !yield(k, v) {
				return
			}
		}
	}
}

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
