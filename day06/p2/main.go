package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common"
	"github.com/mbordner/aoc2018/common/file"
	"strconv"
	"strings"
)

type Pos struct {
	Y int
	X int
}

func (p Pos) String() string {
	return fmt.Sprintf("{%d,%d}", p.X, p.Y)
}

func (p Pos) Add(o Pos) Pos {
	return Pos{Y: p.Y + o.Y, X: p.X + o.X}
}

func (p Pos) Dis(o Pos) int {
	return abs(o.Y-p.Y) + abs(o.X-p.X)
}

func (p Pos) Neighbors() []Pos {
	ns := make([]Pos, 0, 4)
	for _, d := range []Pos{{-1, 0}, {0, 1}, {1, 0}, {0, -1}} {
		ns = append(ns, p.Add(d))
	}
	return ns
}

func abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

func atoi(s string) int {
	val, _ := strconv.ParseInt(s, 10, 64)
	return int(val)
}

type Grid map[Pos]bool

func (g Grid) Has(pos Pos) bool {
	if b, e := g[pos]; e {
		return b
	}
	return false
}

func (g Grid) Seen(pos Pos) bool {
	if _, e := g[pos]; e {
		return true
	}
	return false
}

func dSum(p Pos, ps []Pos) int {
	d := 0
	for _, o := range ps {
		d += p.Dis(o)
	}
	return d
}

func main() {
	ps, maxDistance, pair, center := getData("../data.txt")

	grid := make(Grid)

	queue := make(common.Queue[Pos], 0, 100)

	queue.Enqueue(center)

	area := 0

	for !queue.Empty() {
		cur := *(queue.Dequeue())
		for _, np := range cur.Neighbors() {
			if !grid.Seen(np) {
				if dSum(np, ps) < 10000 {
					grid[np] = true
					area++
					queue.Enqueue(np)
				} else {
					grid[np] = false
				}
			}
		}

	}

	fmt.Println(maxDistance, pair, center)
	fmt.Println(area)

}

func getData(filename string) ([]Pos, int, []Pos, Pos) {
	lines, _ := file.GetLines(filename)
	ps := make([]Pos, len(lines))
	for i, line := range lines {
		tokens := strings.Split(strings.Join(strings.Fields(line), ""), ",")
		ps[i] = Pos{X: atoi(tokens[0]), Y: atoi(tokens[1])}
	}

	pairs := common.GetPairSets(ps)
	maxDist := 0
	var maxPair []Pos
	for _, pair := range pairs {
		d := pair[0].Dis(pair[1])
		if d > maxDist {
			maxDist = d
			maxPair = pair
		}
	}

	cX := 0
	cY := 0

	for _, p := range ps {
		cX += p.X
		cY += p.Y
	}

	cX /= len(ps)
	cY /= len(ps)

	center := Pos{X: cX, Y: cY}

	return ps, maxDist, maxPair, center
}
