package bigexpression

import "math/big"

type variable struct {
	name string
}

func (v variable) Eval(vars map[string]*big.Int) *big.Int {
	if val, exists := vars[v.name]; exists {
		return val
	}
	return big.NewInt(int64(0))
}
