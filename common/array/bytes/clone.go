package bytes

func Clone(a []byte) []byte {
	n := make([]byte, len(a), len(a))
	for i := range a {
		n[i] = a[i]
	}
	return n
}

func Clone2D(a [][]byte) [][]byte {
	n := make([][]byte, len(a), len(a))
	for i := range a {
		n[i] = Clone(a[i])
	}
	return n
}
