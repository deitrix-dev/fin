package iterx

import "iter"

func CollectErr[T any](seq iter.Seq2[T, error]) ([]T, error) {
	var result []T
	for v, err := range seq {
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}

func CollectN[T any](seq iter.Seq[T], n int) []T {
	var result []T
	for v := range seq {
		result = append(result, v)
		if len(result) == n {
			break
		}
	}
	return result
}

func Skip[T any](seq iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		i := 0
		for v := range seq {
			if i >= n {
				if !yield(v) {
					break
				}
			}
			i++
		}
	}
}

func CollectNFilter[T any](seq iter.Seq[T], n int, filter func(T) bool) []T {
	var result []T
	var itersSinceLastFilter int
	for v := range seq {
		if filter(v) {
			result = append(result, v)
			if len(result) == n {
				break
			}
			itersSinceLastFilter = 0
		}
		itersSinceLastFilter++
		if itersSinceLastFilter > 10000 {
			// most of the time, payment iterators extend indefinitely, so we need to break out of
			// the loop if we've iterated over 10000 payments without finding a match. A bit of a
			// hack, but it's a pragmatic solution to a problem that doesn't have a clean solution.
			break
		}
	}
	return result
}

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
