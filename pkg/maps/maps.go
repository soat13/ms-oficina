package maps

func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func Map[T, V any](slice []T, fn func(T) V) []V {
	result := make([]V, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

func MapPtr[T, V any](slice []T, fn func(*T) V) []V {
	result := make([]V, len(slice))
	for i := range slice {
		result[i] = fn(&slice[i])
	}
	return result
}
