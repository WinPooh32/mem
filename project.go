package mem

import "iter"

// Project performs users's projection of the a sequence.
func Project[K, V, Ko, Vo any](a iter.Seq2[K, V], prj func(K, V) (Ko, Vo)) iter.Seq2[Ko, Vo] {
	if prj == nil {
		panic("prj is nil!")
	}
	return func(yield func(k Ko, v Vo) bool) {
		for k, v := range a {
			if !yield(prj(k, v)) {
				return
			}
		}
	}
}
