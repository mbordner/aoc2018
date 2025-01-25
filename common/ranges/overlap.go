package ranges

import (
	"sort"
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

	type OverlapRange[T Number] struct {
		a, b T
	}

	rc := &Collection[T]{}

	for i := 0; i < len(rs)-1; i++ {
		for j := i + 1; j < len(rs); j++ {
			pair := [][]T{rs[i], rs[j]}

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
	}

	return overlaps
}

// Overlaps2 expects an even len array of multiple ranges (len 4 or more), and returns
// an even length array of overlapping ranges with start and end index being the inclusive range boundary
func Overlaps2[T Number](ranges []T) []T {
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

	type OverlapRange[T Number] struct {
		a, b T
	}

	overlapMap := make(map[OverlapRange[T]]bool) // use a map so that we don't generate duplicates
	
	for i := 0; i < len(rs)-1; i++ {
		for j := i + 1; j < len(rs); j++ {
			pair := [][]T{rs[i], rs[j]}

			// ensure pair[0][0] <= pair[1][0]
			if pair[1][0] < pair[0][0] {
				pair[0], pair[1] = pair[1], pair[0]
			}

			if pair[0][0] <= pair[1][0] && pair[0][1] >= pair[1][1] {
				// 2nd pair is contained in first
				overlapMap[OverlapRange[T]{a: Max(pair[0][0], pair[1][0]), b: Min(pair[0][1], pair[1][1])}] = true
			} else if pair[0][1] > pair[1][0] && pair[0][1] <= pair[1][1] {
				overlapMap[OverlapRange[T]{a: pair[1][0], b: Min(pair[0][1], pair[1][1])}] = true
			}
		}
	}

	if len(overlapMap) > 0 { // recurse to check if our overlaps overlap, and resolve this
		overlapRanges := make([]OverlapRange[T], 0, len(overlapMap))

		for overlap := range overlapMap {
			overlapRanges = append(overlapRanges, overlap)
		}

		sort.Slice(overlapRanges, func(i, j int) bool {
			if overlapRanges[i].a < overlapRanges[j].a {
				return true
			} else if overlapRanges[i].a == overlapRanges[j].a {
				return overlapRanges[i].b < overlapRanges[j].b
			}
			return false
		})

		for _, overlap := range overlapRanges {
			overlaps = append(overlaps, []T{overlap.a, overlap.b}...)
		}
	}

	if overlaps != nil && !ArrayEqual(overlaps, ranges) {
		tmpOverlaps := Overlaps[T](overlaps)
		for tmpOverlaps != nil && !ArrayEqual(tmpOverlaps, overlaps) {
			overlaps = tmpOverlaps
			tmpOverlaps = Overlaps[T](overlaps)
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
