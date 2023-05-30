package utils

func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := []K{}
	for k := range m {
		r = append(r, k)
	}
	return r
}
