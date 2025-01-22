package cmath

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Sums(t *testing.T) {
	tests := []struct {
		n        int
		expected []int
	}{
		{
			n:        5,
			expected: []int{0, 5, 1, 4, 2, 3},
		},
		{
			n:        6,
			expected: []int{0, 6, 1, 5, 2, 4, 3, 3},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%d", test.n), func(t *testing.T) {
			vs := Sums(test.n)
			assert.Equal(t, aStr(test.expected), aStr(vs))
		})
	}
}

func aStr(vs []int) string {
	data, _ := json.Marshal(vs)
	return string(data)
}
