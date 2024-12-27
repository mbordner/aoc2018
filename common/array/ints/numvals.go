package ints

import (
	"regexp"
	"strconv"
)

var (
	reDigits = regexp.MustCompile(`\d+`)
)

func NumVals(svals string) []int64 {
	tokens := reDigits.FindAllString(svals, -1)
	vals := make([]int64, len(tokens), len(tokens))
	for i, t := range tokens {
		v, _ := strconv.ParseInt(t, 10, 64)
		vals[i] = v
	}
	return vals
}
