package ints

func Contains(a []int, v int) bool {
	for i := range a {
		if a[i] == v {
			return true
		}
	}
	return false
}
