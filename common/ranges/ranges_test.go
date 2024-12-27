package ranges

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Add(t *testing.T) {

	testCases := []struct {
		name           string
		pairs          [][2]int
		expectedValues [][]int
		expectedLength int
	}{
		{
			name: "initial add",
			pairs: [][2]int{
				{0, 10},
			},
			expectedValues: [][]int{
				{0, 10},
			},
			expectedLength: 11,
		},
		{
			name: "test existing overlaps new",
			pairs: [][2]int{
				{0, 10},
				{2, 8},
			},
			expectedValues: [][]int{
				{0, 10},
				{0, 10},
			},
			expectedLength: 11,
		},
		{
			name: "test extending existing to left and right (overlaps all)",
			pairs: [][2]int{
				{0, 10},
				{2, 8},
				{-5, 15},
			},
			expectedValues: [][]int{
				{0, 10},
				{0, 10},
				{-5, 15},
			},
			expectedLength: 21,
		},
		{
			name: "testing extending to right",
			pairs: [][2]int{
				{0, 10},
				{2, 8},
				{5, 15},
			},
			expectedValues: [][]int{
				{0, 10},
				{0, 10},
				{0, 15},
			},
			expectedLength: 16,
		},
		{
			name: "test extending to left",
			pairs: [][2]int{
				{0, 10},
				{2, 8},
				{-5, 5},
			},
			expectedValues: [][]int{
				{0, 10},
				{0, 10},
				{-5, 10},
			},
			expectedLength: 16,
		},
		{
			name: "test adding to right",
			pairs: [][2]int{
				{0, 10},
				{15, 20},
			},
			expectedValues: [][]int{
				{0, 10},
				{0, 10, 15, 20},
			},
			expectedLength: 17,
		},
		{
			name: "test adding to left",
			pairs: [][2]int{
				{0, 10},
				{-10, -5},
			},
			expectedValues: [][]int{
				{0, 10},
				{-10, -5, 0, 10},
			},
			expectedLength: 17,
		},
		{
			name: "test joining from the middle of two",
			pairs: [][2]int{
				{0, 10},
				{15, 20},
				{5, 17},
			},
			expectedValues: [][]int{
				{0, 10},
				{0, 10, 15, 20},
				{0, 20},
			},
			expectedLength: 21,
		},
		{
			name: "test length minus one",
			pairs: [][2]int{
				{0, 10},
				{12, 20},
			},
			expectedValues: [][]int{
				{0, 10},
				{0, 10, 12, 20},
			},
			expectedLength: 20,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			rc := &Collection[int]{}

			for j, p := range tc.pairs {
				values, err := rc.Add(p[0], p[1])
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedValues[j], values)
			}

			assert.Equal(t, tc.expectedLength, rc.Len())

		})
	}

}
