package util

func Map[T, S any](ts []T, f func(t T) S) []S {
	ss := make([]S, len(ts))

	for i, t := range ts {
		ss[i] = f(t)
	}

	return ss
}
