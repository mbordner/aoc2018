package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestPopulateStringCombinationsAtLength(t *testing.T) {

	tests := []struct {
		chars                string
		uniqueCharsCount     int
		generateLength       int
		expectedCombinations int
	}{
		{
			chars:                "*+",
			generateLength:       2,
			uniqueCharsCount:     2,
			expectedCombinations: 4,
		},
		{
			chars:                "xx",
			generateLength:       2,
			uniqueCharsCount:     1,
			expectedCombinations: 1,
		},
		{
			chars:                "*+",
			generateLength:       3,
			uniqueCharsCount:     2,
			expectedCombinations: 8,
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			results := make(map[string]bool)
			PopulateStringCombinationsAtLength(results, tc.chars, "", tc.generateLength)
			assert.Equal(t, tc.expectedCombinations, len(results))
			assert.Equal(t, int(math.Pow(float64(tc.uniqueCharsCount), float64(tc.generateLength))), len(results))
		})
	}

}

func TestCartesianProduct(t *testing.T) {
	tests := []struct {
		input  [][]string
		output [][]string
	}{
		{
			input:  [][]string{{"a", "b"}, {"c", "d"}},
			output: [][]string{{"a", "c"}, {"a", "d"}, {"b", "c"}, {"b", "d"}},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {

			result := CartesianProduct(tc.input)
			assert.Equal(t, len(tc.output), (result))

		})
	}
}

func TestGetPairs(t *testing.T) {

	tests := []struct {
		values         []int
		expectedLength int
	}{
		{
			values:         []int{1, 2, 3},
			expectedLength: 3,
		},
		{
			values:         []int{1, 2, 3, 4},
			expectedLength: 6,
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			results := GetPairSets[int](tc.values)
			assert.Equal(t, tc.expectedLength, len(results))
			// binomial(n, 2)
			// 1/2 (n - 1) n
			assert.Equal(t, len(tc.values)*(len(tc.values)-1)/2, len(results))
		})
	}

}

func Test_GetRepeatingByteStats(t *testing.T) {

	testCases := []struct {
		s              string
		results        []RepeatingByteStat
		expectedLength int
		min            int
	}{
		{
			s:              "aaabcd",
			results:        []RepeatingByteStat{{b: 'a', i: 0, c: 3}},
			expectedLength: 1,
			min:            3,
		},
		{
			s:              "zdeaaabcd",
			results:        []RepeatingByteStat{{b: 'a', i: 3, c: 3}},
			expectedLength: 1,
			min:            3,
		},
		{
			s:              "aaabcdaaa",
			results:        []RepeatingByteStat{{b: 'a', i: 0, c: 3}, {b: 'a', i: 6, c: 3}},
			expectedLength: 2,
			min:            3,
		},
		{
			s:              "bcdaaa",
			results:        []RepeatingByteStat{{b: 'a', i: 3, c: 3}},
			expectedLength: 1,
			min:            3,
		},
		{
			s:              "bcccccccdaaa",
			results:        []RepeatingByteStat{{b: 'c', i: 1, c: 7}, {b: 'a', i: 9, c: 3}},
			expectedLength: 2,
			min:            3,
		},
		{
			s:              "bcccccccdaaa",
			results:        []RepeatingByteStat{{b: 'c', i: 1, c: 7}},
			expectedLength: 1,
			min:            4,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.s, func(t *testing.T) {

			stats := GetRepeatingByteStats(tc.s, tc.min)

			for i := 0; i < len(stats); i++ {
				assert.Equal(t, tc.results[i], stats[i])
			}

			assert.Equal(t, tc.expectedLength, len(stats))
		})
	}
}
