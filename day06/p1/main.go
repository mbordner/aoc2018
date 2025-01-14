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

type AreaClaim struct {
	id int
	d  int
	p  Pos
}

type AreaClaimGrid map[Pos][]AreaClaim

func (g AreaClaimGrid) CanClaim(id int, p Pos, d int) bool {
	if cs, e := g[p]; e {
		if d > cs[0].d {
			return false
		}
		for _, c := range cs {
			if c.id == id {
				return false
			}
		}
	}
	return true
}

func (g AreaClaimGrid) Claim(id int, p Pos, d int) int {
	if cs, e := g[p]; e {
		if d < cs[0].d {
			g[p] = []AreaClaim{{id: id, d: d, p: p}}
		} else if d == cs[0].d {
			g[p] = append(g[p], AreaClaim{id: id, d: d, p: p})
		}
	} else {
		g[p] = []AreaClaim{{id: id, d: d, p: p}}
	}
	return len(g[p])
}

func main() {
	ps, maxDistance, pair := getData("../data.txt")

	grid := make(AreaClaimGrid)
	for id, p := range ps {
		grid[p] = []AreaClaim{{id: id, d: 0, p: p}}
	}

	queue := make(common.Queue[AreaClaim], 0, 100)
	for _, acs := range grid {
		for _, np := range acs[0].p.Neighbors() {
			queue.Enqueue(AreaClaim{id: acs[0].id, d: 1, p: np})
		}
	}

	maxRange := maxDistance/2 + 2

	for !queue.Empty() {
		curAC := *(queue.Dequeue())
		curD := ps[curAC.id].Dis(curAC.p)
		if curD <= maxRange {
			for _, np := range curAC.p.Neighbors() {
				nAC := AreaClaim{id: curAC.id, d: curAC.d + 1, p: np}
				if grid.CanClaim(nAC.id, nAC.p, nAC.d) {
					grid.Claim(nAC.id, nAC.p, nAC.d)
					queue.Enqueue(nAC)
				}
			}
		}
	}

	finites := make(map[int]bool)
	infinites := make(map[int]bool)

	for _, acs := range grid {
		if len(acs) == 1 {
			if acs[0].d == maxRange {
				for _, ac := range acs {
					infinites[ac.id] = true
				}
			}
		}
	}

	for id := range ps {
		if _, e := infinites[id]; !e {
			finites[id] = true
		}
	}

	areaCounts := make(map[int]int)

	for _, acs := range grid {
		if len(acs) == 1 {
			id := acs[0].id
			if _, finite := finites[id]; finite {
				if c, e := areaCounts[id]; e {
					areaCounts[id] = c + 1
				} else {
					areaCounts[id] = 1
				}
			}
		}
	}

	maxAreaSize := 0
	maxAreaID := 0

	for id, count := range areaCounts {
		if count > maxAreaSize {
			maxAreaSize = count
			maxAreaID = id
		}
	}

	fmt.Println(maxDistance, pair)
	fmt.Println(maxAreaID, maxAreaSize)
}

func getData(filename string) ([]Pos, int, []Pos) {
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

	return ps, maxDist, maxPair
}
