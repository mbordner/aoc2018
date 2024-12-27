package ints

func Remove(a []int, s int) []int {
	r := make([]int, 0, len(a))
	for i := range a {
		if s != a[i] {
			r = append(r, a[i])
		}
	}
	return r
}
