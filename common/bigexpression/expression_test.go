package bigexpression

import (
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func Test_Expressions(t *testing.T) {

	tests := []struct {
		name     string
		expr     string
		input    map[string]*big.Int
		expected int64
	}{
		{
			name:     "addition",
			expr:     "1 + 1",
			input:    make(map[string]*big.Int),
			expected: int64(2),
		},
		{
			name:     "multiplication",
			expr:     "1 * 1",
			input:    make(map[string]*big.Int),
			expected: int64(1),
		},
		{
			name:     "division",
			expr:     "1 / 1",
			input:    make(map[string]*big.Int),
			expected: int64(1),
		},
		{
			name:     "subtraction",
			expr:     "1 - 1",
			input:    make(map[string]*big.Int),
			expected: int64(0),
		},
		{
			name:     "variable addition with constant",
			expr:     "var1 + 3",
			input:    map[string]*big.Int{"var1": big.NewInt(int64(3))},
			expected: int64(6),
		},
		{
			name:     "variable addition with variable",
			expr:     "var1+var2",
			input:    map[string]*big.Int{"var1": big.NewInt(int64(3)), "var2": big.NewInt(int64(4))},
			expected: int64(7),
		},
		{
			name:     "variable multiplication with constant",
			expr:     "var1 * 3",
			input:    map[string]*big.Int{"var1": big.NewInt(int64(3))},
			expected: int64(9),
		},
		{
			name:     "variable multiplication with variable",
			expr:     "var1*var2",
			input:    map[string]*big.Int{"var1": big.NewInt(int64(3)), "var2": big.NewInt(int64(4))},
			expected: int64(12),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			p, err := NewParser(tc.expr)
			assert.Nil(t, err)
			assert.NotNil(t, p)

			v := p.Eval(tc.input)

			assert.Equal(t, tc.expected, v.Int64())

		})
	}

}
