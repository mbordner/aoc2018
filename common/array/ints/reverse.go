package ints

func Reverse(a []int) []int {
	b := Clone(a)
	lh := len(a) / 2
	for i, j := 0, len(b)-1; i < lh; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}
