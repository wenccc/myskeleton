package util

type ArrayItem interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~string | ~float32 | ~float64 | ~bool //这里还有各种 指针型 pointer、通道型 channel、接口型 interface、数组型 array
}

func arrayColumn[T any, V any](s []T, callback func(val T) (target V, find bool)) []V {
	res := make([]V, 0)
	for _, v := range s {
		target, find := callback(v)
		if find {
			res = append(res, target)
		}
	}
	return res
}

func ArrayColumnFind[T any, V int](s []T, callback func(val T) (target V, find bool)) V {
	var res V
	for _, v := range s {
		target, find := callback(v)
		if find {
			if res < target {
				res = target
			}
		}
	}
	return res
}

func InArray[T ArrayItem](needle T, container []T) bool {
	for _, str := range container {
		if str == needle {
			return true
		}
	}
	return false
}
