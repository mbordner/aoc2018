package bytes

type Direction int

const (
	Horizontal Direction = iota
	Vertical
)

func Rotate(a [][]byte) [][]byte {
	n := Clone2D(a)
	for j := 0; j < len(a); j++ {
		for i, c := 0, len(n[j])-1-j; i < len(a); i++ {
			n[i][c] = a[j][i]
		}
	}
	return n
}

func Flip(dir Direction, a [][]byte) (b [][]byte) {
	if dir == Horizontal {
		b = make([][]byte, len(a), len(a))
		for i := range a {
			b[i] = Reverse(a[i])
		}
	} else {
		b = Clone2D(a)
		lh := len(a) / 2
		for i := 0; i < len(a[0]); i++ {
			for j, k := 0, len(a)-1; j < lh; j, k = j+1, k-1 {
				b[j][i], b[k][i] = b[k][i], b[j][i]
			}
		}
	}
	return
}
