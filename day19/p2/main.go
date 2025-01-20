package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common/file"
	"regexp"
	"strconv"
)

var (
	reIpBind    = regexp.MustCompile(`^\#ip\s+(\d+)`)
	reStatement = regexp.MustCompile(`^(\w+)\s+(\d+)\s+(\d+)\s+(\d+)`)
)

type Program []string
type Computer struct {
	regs    []int
	ptr     *int
	program Program
}

func NewComputer(numRegs int) *Computer {
	c := new(Computer)
	c.regs = make([]int, numRegs)
	c.Reset()
	return c
}

func (c *Computer) Reset() {
	for i := range c.regs {
		c.regs[i] = 0
	}
	ptr := 0
	c.ptr = &ptr
	c.program = Program{}
}

func (c *Computer) Load(filename string) {
	c.Reset()
	lines, _ := file.GetLines(filename)
	program := make(Program, 0, len(lines))
	for _, line := range lines {
		if reIpBind.MatchString(line) {
			matches := reIpBind.FindStringSubmatch(line)
			index := atoi(matches[1])
			c.ptr = &(c.regs[index])
		} else if reStatement.MatchString(line) {
			program = append(program, line)
		} else {
			panic("invalid program")
		}
	}
	c.program = program
}

func (c *Computer) Run() {
	for *c.ptr >= 0 && *c.ptr < len(c.program) {
		fmt.Printf("ip=%d [%v] %s ", *c.ptr, c.regs, c.program[*c.ptr])
		matches := reStatement.FindStringSubmatch(c.program[*c.ptr])
		instr, A, B, C := matches[1], atoi(matches[2]), atoi(matches[3]), atoi(matches[4])
		c.EvalInstr(instr, A, B, C)
		fmt.Printf("[%v]\n", c.regs)
		*c.ptr++
	}
}

func (c *Computer) EvalInstr(instr string, A, B, C int) {
	switch instr {
	case "addr": // C = reg(A) + reg(B)
		c.SetRegVal(C, c.GetRegVal(A)+c.GetRegVal(B))
	case "addi": // add reg with immediate
		c.SetRegVal(C, c.GetRegVal(A)+B)
	case "mulr": // mul two regs
		c.SetRegVal(C, c.GetRegVal(A)*c.GetRegVal(B))
	case "muli": // mul reg with immediate
		c.SetRegVal(C, c.GetRegVal(A)*B)
	case "banr": // bitwise and two regs
		c.SetRegVal(C, c.GetRegVal(A)&c.GetRegVal(B))
	case "bani": // bitwise and reg and immediate
		c.SetRegVal(C, c.GetRegVal(A)&B)
	case "borr": // bitwise or two regs
		c.SetRegVal(C, c.GetRegVal(A)|c.GetRegVal(B))
	case "bori": // bitwise or reg and immediate
		c.SetRegVal(C, c.GetRegVal(A)|B)
	case "setr": // set reg value (reg(A) -> C), ignore B
		c.SetRegVal(C, c.GetRegVal(A))
	case "seti": // set reg value (A -> C), ignore b
		c.SetRegVal(C, A)
	case "gtir": // C = A > reg(B) ? 1 : 0
		if A > c.GetRegVal(B) {
			c.SetRegVal(C, 1)
		} else {
			c.SetRegVal(C, 0)
		}
	case "gtri": // C = reg(A) > B ? 1 : 0
		if c.GetRegVal(A) > B {
			c.SetRegVal(C, 1)
		} else {
			c.SetRegVal(C, 0)
		}
	case "gtrr": // C = reg(A) > reg(B) ? 1 : 0
		if c.GetRegVal(A) > c.GetRegVal(B) {
			c.SetRegVal(C, 1)
		} else {
			c.SetRegVal(C, 0)
		}
	case "eqir": // C = A == reg(B) ? 1 : 0
		if A == c.GetRegVal(B) {
			c.SetRegVal(C, 1)
		} else {
			c.SetRegVal(C, 0)
		}
	case "eqri": // C = reg(A) == B ? 1 : 0
		if c.GetRegVal(A) == B {
			c.SetRegVal(C, 1)
		} else {
			c.SetRegVal(C, 0)
		}
	case "eqrr": // C = reg(A) == reg(B) ? 1 : 0
		if c.GetRegVal(A) == c.GetRegVal(B) {
			c.SetRegVal(C, 1)
		} else {
			c.SetRegVal(C, 0)
		}
	default:
		panic("invalid instruction")
	}
}

func (c *Computer) GetRegVal(r int) int {
	if r >= 0 && r < len(c.regs) {
		return c.regs[r]
	}
	panic("invalid reg")
}

func (c *Computer) SetRegVal(r int, val int) {
	if r >= 0 && r < len(c.regs) {
		c.regs[r] = val
	} else {
		panic("invalid reg")
	}
}

func main() {
	c := NewComputer(6)

	c.Load("../data.txt")
	c.SetRegVal(0, 1)

	c.Run()

	fmt.Println(c.GetRegVal(0))

}

func atoi(s string) int {
	val, _ := strconv.ParseInt(s, 10, 64)
	return int(val)
}
