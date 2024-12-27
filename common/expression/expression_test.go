package expression

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Expressions(t *testing.T) {

	tests := []struct {
		name     string
		expr     string
		input    map[string]int64
		expected int64
	}{
		{
			name:     "addition",
			expr:     "1 + 1",
			input:    make(map[string]int64),
			expected: int64(2),
		},
		{
			name:     "multiplication",
			expr:     "1 * 1",
			input:    make(map[string]int64),
			expected: int64(1),
		},
		{
			name:     "division",
			expr:     "1 / 1",
			input:    make(map[string]int64),
			expected: int64(1),
		},
		{
			name:     "subtraction",
			expr:     "1 - 1",
			input:    make(map[string]int64),
			expected: int64(0),
		},
		{
			name:     "Variable addition with constant",
			expr:     "var1 + 3",
			input:    map[string]int64{"var1": int64(3)},
			expected: int64(6),
		},
		{
			name:     "Variable addition with Variable",
			expr:     "var1+var2",
			input:    map[string]int64{"var1": int64(3), "var2": int64(4)},
			expected: int64(7),
		},
		{
			name:     "Variable multiplication with constant",
			expr:     "var1 * 3",
			input:    map[string]int64{"var1": int64(3)},
			expected: int64(9),
		},
		{
			name:     "Variable multiplication with Variable",
			expr:     "var1*var2",
			input:    map[string]int64{"var1": int64(3), "var2": int64(4)},
			expected: int64(12),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			p, err := NewParser(tc.expr)
			assert.Nil(t, err)
			assert.NotNil(t, p)

			v := p.Eval(tc.input)

			assert.Equal(t, tc.expected, v)

		})
	}

}
