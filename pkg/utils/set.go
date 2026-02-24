package utils

import "cmp"

type Set[T cmp.Ordered] map[T]struct{}

func (s Set[T]) Contains(t T) bool {
	_, ok := s[t]
	return ok
}

func (s Set[T]) Add(t T) {
	s[t] = struct{}{}
}

func (s Set[T]) Remove(t T) {
	delete(s, t)
}

func (s Set[T]) Clear() {
	for k := range s {
		delete(s, k)
	}
}

func (s Set[T]) Len() int {
	return len(s)
}
