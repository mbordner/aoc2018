package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common/file"
	"regexp"
	"strconv"
	"strings"
)

var (
	reTest = regexp.MustCompile(`Before:\s+\[((?:\d|,|\s)+)\]\s+((?:\d|\s)+)\s+After:\s+\[((?:\d|,|\s)+)\]`)
)

var (
	instructions = []string{`addr`, `addi`, `mulr`, `muli`, `banr`, `bani`, `borr`, `bori`, `setr`, `seti`, `gtir`, `gtri`, `gtrr`, `eqir`, `eqri`, `eqrr`}
)

type Test struct {
	before    []int
	after     []int
	statement []int
}

type Program [][]int

type OpCodeMap map[int]string

type Computer struct {
	regs    []int
	ptr     int
	program Program
	opCodes OpCodeMap
}

func NewComputer() *Computer {
	c := new(Computer)
	c.Reset()
	return c
}

func (c *Computer) Test(t Test) []string {
	var matched []string
	for _, instr := range instructions {
		c.regs = cloneis(t.before)
		c.EvalInstr(instr, t.statement[1], t.statement[2], t.statement[3])
		m := true
		for i, r := range c.regs {
			if t.after[i] != r {
				m = false
				break
			}
		}
		if m {
			matched = append(matched, instr)
		}
	}
	return matched
}

func (c *Computer) Reset() {
	c.regs = []int{0, 0, 0, 0}
	c.ptr = 0
	c.program = Program{}
}

func (c *Computer) Load(program Program) {
	c.Reset()
	c.program = program
}

func (c *Computer) Run() {
	for c.ptr < len(c.program) {
		op, A, B, C := c.program[c.ptr][0], c.program[c.ptr][1], c.program[c.ptr][2], c.program[c.ptr][3]
		instr := c.opCodes[op]
		c.ptr++
		c.EvalInstr(instr, A, B, C)
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

type InstrPossibility map[string]bool

func (ip InstrPossibility) Eliminate(instr string) {
	if _, e := ip[instr]; e {
		delete(ip, instr)
	}
}

func (ip InstrPossibility) Count() int {
	count := 0
	for _, b := range ip {
		if b {
			count++
		}
	}
	return count
}

func (ip InstrPossibility) Get() string {
	if ip.Count() == 1 {
		for i, b := range ip {
			if b {
				return i
			}
		}
	}
	return ""
}

type OpCodePossibility map[int]InstrPossibility

func (ops OpCodePossibility) Resolved() bool {
	for _, ocps := range ops {
		opInstr := ocps.Get()
		if opInstr == "" {
			return false
		}
	}
	return true
}

func (ops OpCodePossibility) Resolve() {
	for !ops.Resolved() {
		for op, ocps := range ops {
			opInstr := ocps.Get()
			if opInstr != "" {
				for top, tocps := range ops {
					if op != top {
						tocps.Eliminate(opInstr)
					}
				}
			}
		}
	}
}

func main() {
	tests, program := getData("../data.txt")

	c := NewComputer()

	matches := make([][]string, len(tests))
	possibilities := make(OpCodePossibility)
	for i, t := range tests {
		op := t.statement[0]
		matches[i] = c.Test(t)
		if _, e := possibilities[op]; !e {
			possibilities[op] = make(InstrPossibility)
		}
		for _, m := range matches[i] {
			possibilities[op][m] = true
		}
	}

	possibilities.Resolve()

	opCodes := make(OpCodeMap)
	for op, opcp := range possibilities {
		opCodes[op] = opcp.Get()
	}

	if len(opCodes) != len(instructions) {
		panic("not resolved")
	}

	c.opCodes = opCodes
	c.Load(program)

	c.Run()

	fmt.Println(c.GetRegVal(0))

}

func getTest(ts string) Test {
	matches := reTest.FindStringSubmatch(ts)
	before := astois(strings.Split(matches[1], ", "))
	statement := astois(strings.Split(matches[2], " "))
	after := astois(strings.Split(matches[3], ", "))
	return Test{before, after, statement}
}

func getData(filename string) ([]Test, Program) {
	content, _ := file.GetContent(filename)
	var program [][]int

	tokens := strings.Split(strings.TrimSpace(string(content)), "\n\n")
	tests := make([]Test, 0, len(tokens)-1)
	for _, token := range tokens {
		if token == "" {
			continue
		}
		if reTest.MatchString(token) {
			tests = append(tests, getTest(token))
		} else {
			lines := strings.Split(token, "\n")
			program = make([][]int, len(lines))
			for i, line := range lines {
				program[i] = astois(strings.Split(strings.TrimSpace(line), " "))
			}
		}
	}

	return tests, program
}

func atoi(s string) int {
	val, _ := strconv.ParseInt(s, 10, 64)
	return int(val)
}

func astois(ss []string) []int {
	is := make([]int, len(ss))
	for i, s := range ss {
		is[i] = atoi(s)
	}
	return is
}

func cloneis(is []int) []int {
	nis := make([]int, len(is))
	copy(nis, is)
	return nis
}