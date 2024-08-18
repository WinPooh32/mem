package mem

import (
	"iter"

	"golang.org/x/exp/constraints"
)

func Merge[K constraints.Ordered, V any](a, b iter.Seq2[K, V], cs ...iter.Seq2[K, V]) iter.Seq2[K, V] {
	if len(cs) > 0 {
		return mergeN(append(cs, a, b))
	}
	return merge2(a, b)
}

func merge2[K constraints.Ordered, V any](a, b iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(k K, v V) bool) {
		pullA, cancelA := iter.Pull2(a)
		defer cancelA()

		pullB, cancelB := iter.Pull2(b)
		defer cancelB()

		var (
			readyA, readyB         bool
			keyA, keyB, yieldK     K
			valueA, valueB, yieldV V
		)

		for i := 0; ; i++ {
			if !readyA {
				keyA, valueA, readyA = pullA()

				if i > 0 && readyA && keyA < yieldK {
					panic("sequences must be ordered!")
				}
			}
			if !readyB {
				keyB, valueB, readyB = pullB()

				if i > 0 && readyB && keyB < yieldK {
					panic("sequences must be ordered!")
				}
			}

			switch {
			case readyA && readyB:
				if keyA <= keyB {
					readyA = false
					yieldK = keyA
					yieldV = valueA
				} else {
					readyB = false
					yieldK = keyB
					yieldV = valueB
				}
			case readyA:
				readyA = false
				yieldK = keyA
				yieldV = valueA
			case readyB:
				readyB = false
				yieldK = keyB
				yieldV = valueB
			default:
				return
			}

			if !yield(yieldK, yieldV) {
				return
			}
		}
	}
}

func mergeN[K constraints.Ordered, V any](cs []iter.Seq2[K, V]) iter.Seq2[K, V] {
	panic("not implemented!")
}
