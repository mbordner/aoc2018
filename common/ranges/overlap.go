package ranges

import (
	"github.com/mbordner/aoc2018/common"
)

// Overlaps expects an even len array of multiple ranges (len 4 or more), and returns
// an even length array of overlapping ranges where any two input intervals overlapped
func Overlaps[T Number](ranges []T) []T {
	if len(ranges)%2 != 0 { // array length has to be even length, we're expecting pairs [r1_start,r1_end,...
		return []T{}
	}
	if len(ranges) == 2 { // if length is 2, nothing to overlap with
		return ranges
	}

	rs := make([][]T, len(ranges)/2)
	for i, r := 0, 0; i < len(ranges); i, r = i+2, r+1 {
		rs[r] = []T{ranges[i], ranges[i+1]}
		if rs[r][1] < rs[r][0] {
			rs[r][0], rs[r][1] = rs[r][1], rs[r][0]
		}
	}

	var overlaps []T
	pairs := common.GetPairSets(rs)

	type OverlapRange[T Number] struct {
		a, b T
	}

	rc := &Collection[T]{}

	for _, pair := range pairs {
		// ensure pair[0][0] <= pair[1][0]
		if pair[1][0] < pair[0][0] {
			pair[0], pair[1] = pair[1], pair[0]
		}

		if pair[0][0] <= pair[1][0] && pair[0][1] >= pair[1][1] {
			// 2nd pair is contained in first
			overlaps, _ = rc.Add(Max(pair[0][0], pair[1][0]), Min(pair[0][1], pair[1][1]))
		} else if pair[0][1] > pair[1][0] {
			overlaps, _ = rc.Add(pair[1][0], Min(pair[0][1], pair[1][1]))
		}
	}

	return overlaps
}

func Min[T Number](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func Max[T Number](x, y T) T {
	if x > y {
		return x
	}
	return y
}

func ArrayEqual[T Number](x, y []T) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}
