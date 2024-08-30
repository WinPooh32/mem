package mem

import (
	"iter"
)

type Tuple[V1, V2 any] struct {
	L V1
	R []V2
}

// Lookup performs left outer join of a and b sequences using hash algorithm.
// Output order of the a sequence is preserved.
func Lookup[K comparable, V1, V2 any](a iter.Seq2[K, V1], b iter.Seq2[K, V2]) iter.Seq2[K, Tuple[V1, V2]] {
	return func(yield func(k K, v Tuple[V1, V2]) bool) {
		m := map[K][]V2{}
		for k, v := range b {
			m[k] = append(m[k], v)
		}
		for k, v1 := range a {
			tuple := Tuple[V1, V2]{L: v1}
			tuple.R = m[k]
			if !yield(k, tuple) {
				return
			}
		}
	}
}
