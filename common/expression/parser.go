package expression

import (
	"github.com/pkg/errors"
	"regexp"
	"strconv"
)

// https://www.engr.mun.ca/~theo/Misc/exp_parsing.htm

var (
	reSpace       = regexp.MustCompile(`\s`)
	reOperator    = regexp.MustCompile(`\+|\*|\-|\/`)
	reDigits      = regexp.MustCompile(`^\d+$`)
	reVariable    = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9]*$`)
	precedenceMap = map[string]int{
		"*": 10,
		"/": 10,
		"+": 5,
		"-": 5,
	}
)

type Parser struct {
	operators    []*Operator
	operands     []interface{}
	opPrecedence Precedence
	start        int
	end          int
	expr         string
}

func (p *Parser) E() error {
	err := p.P()
	if err != nil {
		return err
	}
	n, err := p.next()
	if err != nil {
		return err
	}
	for reOperator.MatchString(n) && IsBinary(n) {
		p.pushOperator(n)
		p.consume()

		err = p.P()
		if err != nil {
			return err
		}

		n, err = p.next()
		if err != nil {
			return err
		}
	}
	for p.operators[len(p.operators)-1] != nil {
		p.popOperator()
	}
	return nil
}

func (p *Parser) P() error {
	n, err := p.next()
	if err != nil {
		return err
	}
	if reDigits.MatchString(n) {
		v, _ := strconv.ParseInt(n, 10, 64)
		p.operands = append(p.operands, v)
		p.consume()
	} else if reVariable.MatchString(n) {
		p.operands = append(p.operands, Variable{name: n})
		p.consume()
	} else if n == "(" {
		p.consume()
		p.operators = append(p.operators, nil)
		err = p.E()
		if err != nil {
			return err
		}
		n, err = p.next()
		if err != nil {
			return err
		}
		if n != ")" {
			return errors.New("expected )")
		}
		p.consume()
		p.operators = p.operators[0 : len(p.operators)-1]
	} else {
		return errors.New("error parsing expression")
	}
	return nil
}

func (p *Parser) popOperator() {
	op := p.operators[len(p.operators)-1]
	if IsBinary(op.op) {
		p.operators = p.operators[0 : len(p.operators)-1]
		op.right = p.operands[len(p.operands)-1]
		op.left = p.operands[len(p.operands)-2]
		p.operands = p.operands[0 : len(p.operands)-2]
		p.operands = append(p.operands, op)
	}
}

func (p *Parser) pushOperator(op string) {
	for p.operators[len(p.operators)-1] != nil && p.opPrecedence(op, p.operators[len(p.operators)-1].op) <= 0 {
		p.popOperator()
	}
	o := Operator{}
	o.op = op
	p.operators = append(p.operators, &o)
}

func (p *Parser) next() (string, error) {

	for p.start < len(p.expr) && reSpace.MatchString(string(p.expr[p.start])) {
		p.start++
	}

	if p.start == len(p.expr) {
		return "", nil
	}

	p.end = p.start

	if reDigits.MatchString(string(p.expr[p.start])) {
		for p.end < len(p.expr) && reDigits.MatchString(string(p.expr[p.end])) {
			p.end++
		}
	} else if p.expr[p.start] == '(' || p.expr[p.start] == ')' {
		p.end++
	} else if reOperator.MatchString(string(p.expr[p.start])) {
		p.end++
	} else if reVariable.MatchString(string(p.expr[p.start])) {
		for p.end < len(p.expr) && reVariable.MatchString(string(p.expr[p.start:p.end+1])) {
			p.end++
		}
	} else {
		return "", errors.New("unexpected token")
	}

	return p.expr[p.start:p.end], nil
}

func (p *Parser) consume() {
	p.start = p.end
}

func (p *Parser) Eval(vars map[string]int64) int64 {
	return p.operands[0].(*Operator).Eval(vars)
}

func (p *Parser) EvalKnown(vars map[string]int64) (int64, error) {
	return p.operands[0].(*Operator).EvalKnown(vars)
}

func (p *Parser) RootOperator() *Operator {
	return p.operands[0].(*Operator)
}

func (p *Parser) String() string {
	return p.operands[0].(*Operator).String()
}

func NewParser(expr string) (*Parser, error) {
	p := func(op1, op2 string) int {
		if precedenceMap[op1] > precedenceMap[op2] {
			return 1
		}
		if precedenceMap[op1] < precedenceMap[op2] {
			return -1
		}
		return 0
	}
	return NewParserWithPrecedence(expr, p)
}

func NewParserWithPrecedence(expr string, precedence Precedence) (*Parser, error) {
	p := Parser{}
	p.opPrecedence = precedence
	p.expr = expr
	p.operators = make([]*Operator, 0, 20)
	p.operands = make([]interface{}, 0, 20)

	p.operators = append(p.operators, nil)
	err := p.E()
	if err != nil {
		return nil, err
	}

	return &p, nil
}
