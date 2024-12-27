package cmath

type Number interface {
	int | int32 | int64 | float32 | float64
}

func Factorial[V Number](v V) V {
	if v == V(1) {
		return v
	}
	return v * Factorial[V](v-V(1))
}

var (
	MaxInt   = int(^uint(0) >> 1)
	MaxInt64 = int64(^uint64(0) >> 1)
)
