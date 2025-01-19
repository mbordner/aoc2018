package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common"
	"github.com/mbordner/aoc2018/common/file"
	"regexp"
	"strconv"
)

var (
	reClayLoc = regexp.MustCompile(`(x|y)=(\d+),\s+(x|y)=(\d+)\.{2}(\d+)`)
)

const (
	Sand      = '.'
	Clay      = '#'
	Spring    = '+'
	WaterRest = '~'
	WaterFlow = '|'
)

func main() {
	spring, springGrid, grid := getData("../data.txt")

	flow(grid, common.Pos{X: spring.X, Y: 0})

	fmt.Println(springGrid)
	fmt.Println(grid)

	count := 0
	for y := range grid {
		for x := range grid[y] {
			if grid[y][x] == WaterRest {
				count++
			}
		}
	}

	fmt.Printf("Number of water at rest tiles is %d.\n", count)
}

func is(g common.Grid, p common.Pos, vals []byte) bool {
	if !g.ContainsPos(p) {
		return false
	}
	pv := g.Val(p)
	for _, v := range vals {
		if v == pv {
			return true
		}
	}
	return false
}

type ExtendReason int

const (
	YesExtend ExtendReason = iota
	NoHitsClay
	NoOverFlowsDown
	NoOffGrid
)

func canRowExtend(g common.Grid, p common.Pos) ExtendReason {
	if !g.ContainsPos(p) {
		return NoOffGrid
	}
	bp := p.Add(common.DD)
	if g.ContainsPos(bp) && is(g, bp, []byte(string(Sand))) {
		return NoOverFlowsDown
	}
	if is(g, p, []byte{Clay}) {
		return NoHitsClay
	}
	return YesExtend
}

func flow(g common.Grid, p common.Pos) {
	g.Set(p, WaterFlow)
	bp := p.Add(common.DD)

	if is(g, bp, []byte{Sand}) {
		flow(g, bp)
	}

	if is(g, bp, []byte{WaterRest, Clay}) {
		var left, right common.Pos
		leftExtendReason, rightExtendReason := YesExtend, YesExtend
		// extend left, stop if left is clay, or below is sand
		left = p
		for leftExtendReason == YesExtend {
			if leftExtendReason = canRowExtend(g, left); leftExtendReason == YesExtend {
				left = left.Add(common.DL)
			}
		}
		// extend right, stop if right is clay, or below is sand
		right = p
		for rightExtendReason == YesExtend {
			if rightExtendReason = canRowExtend(g, right); rightExtendReason == YesExtend {
				right = right.Add(common.DR)
			}
		}

		if leftExtendReason == NoHitsClay && rightExtendReason == NoHitsClay {
			for np := p; np != left; np = np.Add(common.DL) {
				g.Set(np, WaterRest)
			}
			for np := p; np != right; np = np.Add(common.DR) {
				g.Set(np, WaterRest)
			}
		} else {
			for np := left.Add(common.DR); np != right; np = np.Add(common.DR) {
				g.Set(np, WaterFlow)
			}
			if leftExtendReason == NoOverFlowsDown {
				flow(g, left)
			}
			if rightExtendReason == NoOverFlowsDown {
				flow(g, right)
			}
		}

	}
}

func getData(filename string) (common.Pos, common.Grid, common.Grid) {
	pc := make(common.PosContainer)

	lines, _ := file.GetLines(filename)
	for _, line := range lines {
		matches := reClayLoc.FindStringSubmatch(line)
		a := atoi(matches[2])
		b1 := atoi(matches[4])
		b2 := atoi(matches[5])
		for b := b1; b <= b2; b++ {
			if matches[1] == "x" && matches[3] == "y" {
				pc[common.Pos{Y: b, X: a}] = true
			} else if matches[1] == "y" && matches[3] == "x" {
				pc[common.Pos{Y: a, X: b}] = true
			}
		}
	}

	spring := common.Pos{Y: 0, X: 500}

	minP, maxP := pc.Extents()

	minP.X-- // add space to left for overflow, will be sand
	maxP.X++ // add space to right for overflow, will be sand

	w := maxP.X - minP.X + 1
	h := maxP.Y - minP.Y + 1

	grid := make(common.Grid, h)
	for y := range grid {
		grid[y] = make([]byte, w)
		for x := range grid[y] {
			grid[y][x] = Sand
		}
	}

	spring = spring.Sub(minP)

	for p := range pc {
		p = p.Sub(minP)
		grid[p.Y][p.X] = Clay
	}

	spc := make(common.PosContainer)
	spc[common.Pos{Y: spring.Y, X: minP.X}] = true
	spc[common.Pos{Y: minP.Y - 1, X: maxP.X}] = true
	tSpring := common.Pos{Y: 0, X: 500}
	minSP, maxSP := spc.Extents()
	sw := maxSP.X - minSP.X + 1
	sh := maxSP.Y - minSP.Y + 1
	springGrid := make(common.Grid, sh)
	springX := tSpring.X - minSP.X
	for y := range springGrid {
		springGrid[y] = make([]byte, sw)
		for x := range springGrid[y] {
			springGrid[y][x] = Sand
		}
		springGrid[y][springX] = Spring
	}

	return spring, springGrid, grid
}

func atoi(s string) int {
	val, _ := strconv.ParseInt(s, 10, 64)
	return int(val)
}
