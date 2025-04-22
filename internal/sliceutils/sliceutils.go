package sliceutils

import "iter"

type pair[E any] struct {
	One E
	Two E
}

func Pairs[E any](s []E) iter.Seq2[int, pair[E]] {
	return func(yield func(int, pair[E]) bool) {
		for i := 0; (i*2)+1 < len(s); i++ {
			p := pair[E]{s[i*2], s[i*2+1]}
			if !yield(i, p) {
				return
			}
		}
	}
}
