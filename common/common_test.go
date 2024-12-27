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
