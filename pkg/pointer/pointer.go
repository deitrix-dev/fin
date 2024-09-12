package pointer

func To[T any](v T) *T {
	return &v
}

func Zero[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}
