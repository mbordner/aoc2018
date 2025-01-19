package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common"
	"github.com/mbordner/aoc2018/common/file"
)

const (
	Open       = '.'
	Trees      = '|'
	Lumberyard = '#'
)

type LandCounts map[byte]int

func (lc LandCounts) Count(char byte) int {
	if c, e := lc[char]; e {
		return c
	}
	return 0
}

func main() {
	acres := getData("../data.txt")
	fmt.Printf("Initial State:\n")
	fmt.Println(acres)

	mins := 0
	for mins < 10 {
		acres.Tick()
		mins++
		fmt.Printf("After %d mins\n", mins)
		fmt.Println(acres)
	}

	lc := acres.LandCounts()
	fmt.Printf("Resource value: %d\n", lc.Count(Trees)*lc.Count(Lumberyard))
}

type Acres struct {
	g common.Grid
}

func (a *Acres) LandCounts() LandCounts {
	lc := LandCounts{}
	for y := range a.g {
		for x := range a.g[y] {
			char := a.g[y][x]
			lc[char] = lc.Count(char) + 1
		}
	}
	return lc
}

func (a *Acres) String() string {
	return a.g.String()
}

func (a *Acres) Tick() {
	g := a.g.Clone()

	for y := range a.g {
		for x := range a.g[y] {
			char := a.g[y][x]
			lc := a.getAdjacentCounts(y, x)
			switch char {
			case Open:
				if lc.Count(Trees) >= 3 {
					g[y][x] = Trees
				} else {
					g[y][x] = Open
				}
			case Trees:
				if lc.Count(Lumberyard) >= 3 {
					g[y][x] = Lumberyard
				} else {
					g[y][x] = Trees
				}
			case Lumberyard:
				if lc.Count(Trees) >= 1 && lc.Count(Lumberyard) >= 1 {
					g[y][x] = Lumberyard
				} else {
					g[y][x] = Open
				}
			}
		}
	}

	a.g = g
}

func (a *Acres) getAdjacentCounts(y, x int) LandCounts {
	counts := make(LandCounts)
	p := common.Pos{Y: y, X: x}
	adjacent := p.AdjacentWithCorners()
	for _, ap := range adjacent {
		if a.g.ContainsPos(ap) {
			char := a.g[ap.Y][ap.X]
			counts[char] = counts.Count(char) + 1
		}
	}
	return counts
}

func getData(filename string) *Acres {
	a := Acres{}
	lines, _ := file.GetLines(filename)
	a.g = common.ConvertGrid(lines)
	return &a
}
