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

func CollectNErr[T any](seq iter.Seq2[T, error], n int) ([]T, error) {
	var result []T
	for v, err := range seq {
		if err != nil {
			return nil, err
		}
		result = append(result, v)
		if len(result) == n {
			break
		}
	}
	return result, nil
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

func SkipErr[T any](seq iter.Seq2[T, error], n int) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		i := 0
		for v, err := range seq {
			if err != nil {
				if !yield(v, err) {
					break
				}
			}
			if i >= n {
				if !yield(v, err) {
					break
				}
			}
			i++
		}
	}
}

func Paginate[T any](seq iter.Seq[T], offset, limit uint) []T {
	return CollectN(Skip(seq, int(offset)), int(limit))
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

// Join returns a new iterator that yields the values from the input iterators in sorted order.
// Join assumes that the input iterators are already sorted according to the cmp function.
func Join[T any](cmp func(T, T) int, seqs ...iter.Seq[T]) iter.Seq[T] {
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

func JoinErr[T any](cmp func(T, T) int, seqs ...iter.Seq2[T, error]) iter.Seq2[T, error] {
	type puller struct {
		val  T
		err  error
		ok   bool
		next func() (T, error, bool)
		stop func()
	}

	return func(yield func(T, error) bool) {
		// prime the iterators by pulling the first value from each.
		pullers := make([]puller, len(seqs))
		for i, seq := range seqs {
			next, stop := iter.Pull2(seq)
			v, err, ok := next()
			pullers[i] = puller{val: v, err: err, ok: ok, next: next, stop: stop}
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
			v, err := pullers[minIndex].val, pullers[minIndex].err
			if !yield(v, err) {
				break
			}
			// advance the puller that yielded the smallest value
			p := &pullers[minIndex]
			p.val, p.err, p.ok = p.next()
		}
	}
}

func WithNilErr[T any](seq iter.Seq[T]) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for v := range seq {
			if !yield(v, nil) {
				break
			}
		}
	}
}

// FirstError returns a new iterator that yields the values from the input iterator until an error
// is encountered. If an error is encountered, the error is returned from the iterator's close
// function.
func FirstError[T any](seq iter.Seq2[T, error], err *error) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v, e := range seq {
			if e != nil {
				*err = e
				return
			}
			if !yield(v) {
				return
			}
		}
	}
}
