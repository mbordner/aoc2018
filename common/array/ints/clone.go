package ints

func Clone(a []int) []int {
	n := make([]int, len(a), len(a))
	for i := range a {
		n[i] = a[i]
	}
	return n
}

func Clone2D(a [][]int) [][]int {
	n := make([][]int, len(a), len(a))
	for i := range a {
		n[i] = Clone(a[i])
	}
	return n
}
