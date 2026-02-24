package fns

func Map[T any, R any](sl []T, fn func(T) R) []R {
	n := make([]R, len(sl))
	for i, v := range sl {
		n[i] = fn(v)
	}
	return n
}
