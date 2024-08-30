package mem

import (
	"cmp"
	"iter"
)

// Merge combines non-decreasing sequences into one.
// Will panic when some of sequences are not ordered.
func Merge[K cmp.Ordered, V any](a, b iter.Seq2[K, V], cs ...iter.Seq2[K, V]) iter.Seq2[K, V] {
	if len(cs) > 0 {
		return mergeN(append(cs, a, b))
	}
	return merge2(a, b)
}

func merge2[K cmp.Ordered, V any](a, b iter.Seq2[K, V]) iter.Seq2[K, V] {
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

func mergeN[K cmp.Ordered, V any](cs []iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(k K, v V) bool) {
		pulls := make([]func() (K, V, bool), 0, len(cs))
		cancels := make([]func(), 0, len(cs))

		for _, seq := range cs {
			pull, cancel := iter.Pull2(seq)
			pulls = append(pulls, pull)
			cancels = append(cancels, cancel)
		}

		defer func() {
			for _, cancel := range cancels {
				cancel()
			}
		}()

		var (
			keys  = make([]K, len(cs))
			vals  = make([]V, len(cs))
			ready = make([]bool, len(cs))

			last K
		)

		for i := 0; ; i++ {
			pull(pulls, keys, vals, ready)

			readyMinIdx := argMin(keys, ready)
			if readyMinIdx < 0 {
				return
			}

			curKey := keys[readyMinIdx]
			if curKey < last {
				panic("sequences must be ordered!")
			}

			if !yield(curKey, vals[readyMinIdx]) {
				return
			}

			last = curKey
			ready[readyMinIdx] = false
		}
	}
}

func pull[K cmp.Ordered, V any](pulls []func() (K, V, bool), keys []K, vals []V, ready []bool) {
	if len(pulls) != len(keys) || len(pulls) != len(vals) || len(pulls) != len(ready) {
		panic("all lengths must be equal!")
	}
	for i, pull := range pulls {
		if ready[i] {
			continue
		}
		keys[i], vals[i], ready[i] = pull()
	}
}

func argMin[K cmp.Ordered](x []K, ready []bool) int {
	if len(x) < 1 {
		panic("empty!")
	}
	if len(x) != len(ready) {
		panic("lengths must be equal!")
	}
	var p int
	// Search for the first ready value.
	for p = 0; p < len(ready) && !ready[p]; p++ {
	}
	if p == len(ready) {
		// All values are not ready.
		return -1
	}
	m := x[p]
	for i := p + 1; i < len(x); i++ {
		if ready[i] {
			if v := x[i]; v < m {
				m = v
				p = i
			}
		}
	}
	return p
}
