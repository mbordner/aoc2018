package ranges

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Overlaps(t *testing.T) {

	tests := []struct {
		input    []int
		expected []int
	}{
		{
			input:    []int{-3, 4, 0, 4, -3, 7},
			expected: []int{-3, 4},
		},
		{
			input:    []int{1, 3, 2, 4, 0, 6},
			expected: []int{1, 4},
		},
		{
			input:    []int{0, 4, 8, 9, -3, 5},
			expected: []int{0, 4},
		},
		{
			input:    []int{1, 3, 2, 4},
			expected: []int{2, 3},
		},
		{
			input:    []int{1, 3, 2, 5, 4, 6},
			expected: []int{2, 5},
		},
		{
			input:    []int{10, 20, 15, 40},
			expected: []int{15, 20},
		},
		{
			input:    []int{10, 20, 15, 25, 30, 45, 40, 50},
			expected: []int{15, 20, 40, 45},
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v", tc.input), func(t *testing.T) {
			results := Overlaps[int](tc.input)
			assert.NotNil(t, results)
			assert.True(t, ArrayEqual[int](results, tc.expected))
		})
	}
}
