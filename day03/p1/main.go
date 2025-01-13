package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common/file"
	"regexp"
	"strconv"
)

var (
	reClaim = regexp.MustCompile(`#(\d+) @ (\d+),(\d+): (\d+)x(\d+)`)
)

type Claim struct {
	id int
	x  int
	y  int
	w  int
	h  int
}

type Grid struct {
	data [][][]int
}

func NewGrid(sz int) *Grid {
	grid := Grid{}
	grid.data = make([][][]int, sz)
	for y := range grid.data {
		grid.data[y] = make([][]int, sz)
		for x := range grid.data[y] {
			grid.data[y][x] = []int{}
		}
	}
	return &grid
}

func (g *Grid) Claim(id, cx, cy, cw, ch int) {
	for y := cy; y < cy+ch; y++ {
		for x := cx; x < cx+cw; x++ {
			g.data[y][x] = append(g.data[y][x], id)
		}
	}
}

func (g *Grid) GetOverlapping() int {
	count := 0
	for y := range g.data {
		for x := range g.data[y] {
			if len(g.data[y][x]) > 1 {
				count++
			}
		}
	}
	return count
}

func main() {
	grid := NewGrid(1000)
	claims := getClaims("../data.txt")
	for _, c := range claims {
		grid.Claim(c.id, c.x, c.y, c.w, c.h)
	}

	fmt.Println(grid.GetOverlapping())
}

func getClaims(filename string) []Claim {
	lines, _ := file.GetLines(filename)
	claims := make([]Claim, len(lines))
	for i, line := range lines {
		matches := reClaim.FindStringSubmatch(line)
		claims[i] = Claim{id: atoi(matches[1]), x: atoi(matches[2]), y: atoi(matches[3]), w: atoi(matches[4]), h: atoi(matches[5])}
	}
	return claims
}

func atoi(s string) int {
	val, _ := strconv.ParseInt(s, 10, 64)
	return int(val)
}
