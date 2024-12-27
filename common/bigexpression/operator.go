package bigexpression

import (
	"github.com/pkg/errors"
	"math/big"
)

// Precedence returns > 0 if op1 > op2, or < 0 if op1 < op2, otherwise 0
type Precedence func(op1, op2 string) int

type Operator struct {
	op    string
	left  interface{}
	right interface{}
}

func (o *Operator) Eval(vars map[string]*big.Int) *big.Int {
	var l, r *big.Int
	switch tl := o.left.(type) {
	case variable:
		l = tl.Eval(vars)
	case *big.Int:
		l = tl
	case *Operator:
		l = tl.Eval(vars)
	}
	switch tr := o.right.(type) {
	case variable:
		r = tr.Eval(vars)
	case *big.Int:
		r = tr
	case *Operator:
		r = tr.Eval(vars)
	}
	switch o.op {
	case "-":
		return big.NewInt(0).Sub(l, r)
	case "+":
		return big.NewInt(0).Add(l, r)
	case "*":
		return big.NewInt(0).Mul(l, r)
	case "/":
		return big.NewInt(0).Div(l, r)
	case "|":
		n := new(big.Int)
		n, _ = n.SetString(l.String()+r.String(), 10)
		return n
	}
	panic(errors.New("unknown operator"))
}

func IsBinary(op string) bool {
	return true
}
