package array

import (
	"cmp"
	"slices"
	"strconv"
	"strings"
	"unsafe"
)

type Number interface {
	int | int32 | int64 | float32 | float64
}

func ToNumbers[V Number](s, sep string) []V {
	tokens := strings.Split(s, sep)
	nums := make([]V, len(tokens), len(tokens))
	var t V
	stbits := 8 * int(unsafe.Sizeof(t))
	for i := range tokens {
		val, _ := strconv.ParseInt(tokens[i], 10, stbits)
		nums[i] = V(val)
	}
	return nums
}

func CloneNumbers[V Number](a []V) []V {
	n := make([]V, len(a), len(a))
	copy(n, a)
	return n
}

func ReverseNumbers[V Number](a []V) []V {
	b := CloneNumbers[V](a)
	lh := len(a) / 2
	for i, j := 0, len(b)-1; i < lh; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}

func SumNumbers[V Number](a []V) V {
	s := V(0)
	for _, v := range a {
		s += v
	}
	return s
}

func AllSameNumbers[V Number](a []V) bool {
	for i := 1; i < len(a); i++ {
		if a[i] != a[i-1] {
			return false
		}
	}
	return true
}

func Equals[V Number](a, b []V) bool {
	if len(a) != len(b) {
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
	}
	return true
}

func Clone[A any](a []A) []A {
	n := make([]A, len(a), len(a))
	copy(n, a)
	return n
}

func Reverse[A any](a []A) []A {
	b := Clone[A](a)
	lh := len(a) / 2
	for i, j := 0, len(b)-1; i < lh; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}

func SwapRowCols[A any](a [][]A) [][]A {
	cl := len(a)
	rl := len(a[0])

	b := make([][]A, rl)
	for j := 0; j < rl; j++ {
		b[j] = make([]A, cl)
		for i := 0; i < cl; i++ {
			b[j][i] = a[i][j]
		}
	}

	return b
}

func Keys[A comparable, V any](m map[A]V) []A {
	a := make([]A, len(m))
	i := 0
	for k, _ := range m {
		a[i] = k
		i++
	}
	return a
}

func SortedKeys[A cmp.Ordered, V any](m map[A]V) []A {
	keys := Keys(m)
	slices.Sort(keys)
	return keys
}

func Values[A comparable, V any](m map[A]V) []V {
	a := make([]V, len(m))
	i := 0
	for _, v := range m {
		a[i] = v
		i++
	}
	return a
}

// ( N * N-1 ) / 2
func Pairs[V any](a []V) [][]V {
	b := make([][]V, 0, (len(a)*len(a)-1)/2)
	for j := 0; j < len(a); j++ {
		for i := j + 1; i < len(a); i++ {
			b = append(b, []V{a[j], a[i]})
		}
	}
	return b
}

func Contains[V comparable](h []V, n V) bool {
	for _, v := range h {
		if v == n {
			return true
		}
	}
	return false
}
