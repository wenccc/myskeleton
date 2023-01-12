package util

type MapKey interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~string | ~float32 | ~float64 | ~bool //这里还有各种 指针型 pointer、通道型 channel、接口型 interface、数组型 array
}

func MapSetIfNotExist[T MapKey, V any](m map[T]V, key T, Value V) bool {
	if _, ok := m[key]; ok {
		return true
	}
	m[key] = Value
	return false
}
