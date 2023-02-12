package helpers

func RemoveElementByIndex[T any](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

func FindIndex[T any](slice []T, matchFunc func(T) bool) int {
	for index, element := range slice {
		if matchFunc(element) {
			return index
		}
	}

	return -1
}

func FirstOrDefault[T any](slice []T, filter func(*T) bool) (element *T) {

	for i := 0; i < len(slice); i++ {
		if filter(&slice[i]) {
			return &slice[i]
		}
	}

	return nil
}

func Where[T any](slice []T, filter func(*T) bool) []*T {

	var ret []*T = make([]*T, 0)

	for i := 0; i < len(slice); i++ {
		if filter(&slice[i]) {
			ret = append(ret, &slice[i])
		}
	}

	return ret
}
