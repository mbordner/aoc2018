package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
)

func PopulateStringCombinationsAtLength(results map[string]bool, pickChars string, prefix string, length int) {
	if length == 0 {
		results[prefix] = true
		return
	}

	for i := 0; i < len(pickChars); i++ {
		nextPrefix := prefix + string(pickChars[i])
		PopulateStringCombinationsAtLength(results, pickChars, nextPrefix, length-1)
	}
}

type Anything interface{}

func GetPairSets[T Anything](elements []T) [][]T {
	values := make([][]T, 0, len(elements)*(len(elements)-1)/2)
	for i := 0; i < len(elements)-1; i++ {
		for j := i + 1; j < len(elements); j++ {
			values = append(values, []T{elements[i], elements[j]})
		}
	}
	return values
}

func CartesianProduct[T any](sets [][]T) [][]T {
	result := [][]T{{}}

	for _, set := range sets {
		temp := [][]T{}
		for _, x := range set {
			for _, r := range result {
				temp = append(temp, append(r, x))
			}
		}
		result = temp
	}

	return result
}

func FilterArray[T comparable](array []T, exclude []T) []T {
	ex := make(map[T]bool)
	for _, x := range exclude {
		ex[x] = true
	}

	values := make([]T, 0, len(array))

	for _, v := range array {
		if _, e := ex[v]; !e {
			values = append(values, v)
		}
	}

	return values
}

type IntNumber interface {
	int | int32 | int64
}

func Min[T IntNumber](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T IntNumber](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Abs[T IntNumber](v T) T {
	if v >= 0 {
		return v
	}
	return T(-1) * v
}

type Grid [][]byte
type Pos struct {
	Y int
	X int
}

func (p Pos) Add(o Pos) Pos {
	return Pos{Y: p.Y + o.Y, X: p.X + o.X}
}

func (p Pos) Scale(s int) Pos {
	return Pos{Y: p.Y * s, X: p.X * s}
}

type PosDelta Pos

var (
	DN = Pos{Y: -1, X: 0}
	DU = DN
	DE = Pos{Y: 0, X: 1}
	DR = DE
	DS = Pos{Y: 1, X: 0}
	DD = DS
	DW = Pos{Y: 0, X: -1}
	DL = DW
)

func (p Pos) String() string {
	return fmt.Sprintf("{%d,%d}", p.X, p.Y)
}

type Positions []Pos

func (p Positions) Extents() (Pos, Pos) {
	var min, max Pos = p[0], p[0]
	for i := 1; i < len(p); i++ {
		o := p[i]
		if o.Y < min.Y {
			min = o
		}
		if o.X < min.X {
			min = o
		}
		if o.Y > max.Y {
			max = o
		}
		if o.X > max.X {
			max = o
		}
	}
	return min, max
}

func (g Grid) Contains(x, y int) bool {
	if y >= 0 && y < len(g) && x >= 0 && x < len(g[y]) {
		return true
	}
	return false
}

func (g Grid) ContainsPos(p Pos) bool {
	return g.Contains(p.X, p.Y)
}

func (g Grid) Print() {
	for _, row := range g {
		fmt.Println(string(row))
	}
}

func ConvertGrid(lines []string) Grid {
	grid := make(Grid, len(lines))
	for i, line := range lines {
		grid[i] = []byte(line)
	}
	return grid
}

type Queue[T any] []T

func (q *Queue[T]) Enqueue(s T) {
	*q = append(*q, s)
}

func (q *Queue[T]) Empty() bool {
	return len(*q) == 0
}

func (q *Queue[T]) Dequeue() *T {
	if !q.Empty() {
		s := (*q)[0]
		*q = (*q)[1:]
		return &s
	}
	return nil
}

type Stack[T any] []T

func (s *Stack[T]) Len() int {
	return len(*s)
}

func (s *Stack[T]) Push(v T) {
	*s = append(*s, v)
}

func (s *Stack[T]) Peek() *T {
	if !s.Empty() {
		v := (*s)[len(*s)-1]
		return &v
	}
	return nil
}

func (s *Stack[T]) Empty() bool {
	return len(*s) == 0
}

func (s *Stack[T]) Pop() *T {
	if !s.Empty() {
		v := (*s)[len(*s)-1]
		*s = (*s)[:len(*s)-1]
		return &v
	}
	return nil
}

type PosContainer map[Pos]bool

func (v PosContainer) Has(p Pos) bool {
	if b, e := v[p]; e {
		return b
	}
	return false
}

type PosLinker map[Pos]Pos

func Filter[T comparable](values []T, value T) []T {
	vs := make([]T, 0, len(values))
	for _, v := range values {
		if v != value {
			vs = append(vs, v)
		}
	}
	return vs
}

func Dedupe[T comparable](values []T) []T {
	vm := make(map[T]bool)
	for _, v := range values {
		vm[v] = true
	}
	vs := make([]T, 0, len(values))
	for v := range vm {
		vs = append(vs, v)
	}
	return vs
}

func ByteCharToInt(char byte) int {
	return int(char - '0')
}

func StrToA(str string) int {
	return int(StrToA64(str))
}

func StrToA64(str string) int64 {
	val, _ := strconv.ParseInt(str, 10, 64)
	return val
}

func HashString(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

type RepeatingByteStat struct {
	b byte
	i int
	c int
}

func GetRepeatingByteStats(s string, min int) []RepeatingByteStat {

	rbs := make([]RepeatingByteStat, 0, 10)

	last := s[0]
	count := 0
	for i := 1; i < len(s); i++ {
		if s[i] == last {
			if count == 0 {
				count = 2
			} else {
				count++
			}
		} else if count > 0 {
			if count >= min {
				rbs = append(rbs, RepeatingByteStat{b: last, i: i - count, c: count})
			}
			count = 0
		}
		last = s[i]
	}

	if count >= min {
		rbs = append(rbs, RepeatingByteStat{b: last, i: len(s) - count, c: count})
	}

	return rbs
}
