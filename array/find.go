package array

func Find[T comparable](sources []T, element T) (index int, ok bool) {
	for i, item := range sources {
		if item == element {
			return i, true
		}
	}
	return -1, false
}
