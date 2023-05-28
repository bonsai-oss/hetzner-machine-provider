package helper

func Filter[T any](s []T, fn func(T) bool) []T {
	var p []T
	for _, v := range s {
		if fn(v) {
			p = append(p, v)
		}
	}
	return p
}

func Map[T any, U any](s []T, fn func(T) U) []U {
	var p []U
	for _, v := range s {
		p = append(p, fn(v))
	}
	return p
}
