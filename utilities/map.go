package utilities

type Number interface {
	int | int64 | int32 | uint | uint64 | uint32 | float64 | float32
}

type MapKey interface {
	string | Number
}

func GetKeys[K MapKey, V any](m map[K]V) []K {
	var keys []K
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
