package iterx

import "iter"

func Empty[T any]() iter.Seq[T] {
	return func(yield func(T) bool) {}
}

// JoinFunc returns a new iterator that yields the values from the input iterators in sorted order.
// JoinFunc assumes that the input iterators are already sorted according to the cmp function.
func JoinFunc[T any](seqs []iter.Seq[T], cmp func(T, T) int) iter.Seq[T] {
	type puller struct {
		val  T
		ok   bool
		next func() (T, bool)
		stop func()
	}

	return func(yield func(T) bool) {
		// prime the iterators by pulling the first value from each.
		pullers := make([]puller, len(seqs))
		for i, seq := range seqs {
			next, stop := iter.Pull(seq)
			v, ok := next()
			pullers[i] = puller{val: v, ok: ok, next: next, stop: stop}
		}
		defer func() {
			for _, p := range pullers {
				p.stop()
			}
		}()

		for {
			// find the non-empty puller with the smallest value, according to the cmp function.
			minIndex := -1
			for i, p := range pullers {
				if !p.ok {
					continue
				}
				if minIndex == -1 || cmp(p.val, pullers[minIndex].val) == -1 {
					minIndex = i
				}
			}
			if minIndex == -1 {
				// all pullers are empty
				break
			}
			// send the smallest value to the caller
			v := pullers[minIndex].val
			if !yield(v) {
				break
			}
			// advance the puller that yielded the smallest value
			p := &pullers[minIndex]
			p.val, p.ok = p.next()
		}
	}
}
