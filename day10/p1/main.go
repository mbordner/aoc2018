package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common/file"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	reLight = regexp.MustCompile(`position=<\s*(-?\d+)\s*,\s*(-?\d+)\s*> velocity=<\s*(-?\d+)\s*,\s*(-?\d+)\s*>`)
)

type Vector struct {
	X int
	Y int
}

type Vectors []Vector

func (v Vector) Add(o Vector) Vector {
	return Vector{v.X + o.X, v.Y + o.Y}
}

func (v Vector) Scale(s int) Vector {
	return Vector{v.X * s, v.Y * s}
}

func (v Vector) Dis(o Vector) int {
	return abs(v.X-o.X) + abs(v.Y-o.Y)
}

func (v Vectors) Add(os Vectors) {
	for i := range v {
		v[i] = v[i].Add(os[i])
	}
}

func (v Vectors) Scale(s int) {
	for i := range v {
		v[i] = v[i].Scale(s)
	}
}

func (v Vectors) Extents() (Vector, Vector) {
	minV, maxV := v[0], v[0]
	for i := 1; i < len(v); i++ {
		if v[i].X < minV.X {
			minV.X = v[i].X
		}
		if v[i].X > maxV.X {
			maxV.X = v[i].X
		}
		if v[i].Y < minV.Y {
			minV.Y = v[i].Y
		}
		if v[i].Y > maxV.Y {
			maxV.Y = v[i].Y
		}
	}

	return minV, maxV
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

func NewGrid(ps Vectors, off, on byte) string {
	minV, maxV := ps.Extents()
	return NewGridWithExtents(ps, minV, maxV, off, on)
}

func NewGridWithExtents(ps Vectors, minV Vector, maxV Vector, off, on byte) string {
	w := maxV.X - minV.X + 1
	h := maxV.Y - minV.Y + 1

	grid := make([][]byte, h)
	for y := 0; y < h; y++ {
		grid[y] = make([]byte, w)
		for x := 0; x < w; x++ {
			grid[y][x] = off
		}
	}

	for _, v := range ps {
		y := v.Y - minV.Y
		x := v.X - minV.X
		grid[y][x] = on
	}

	lines := make([]string, len(grid))
	for y, line := range grid {
		lines[y] = string(line)
	}

	return strings.Join(lines, "\n")
}

func main() {
	ps, vs := getData("../data.txt")

	var b []byte = make([]byte, 1)

	secs := 0
	for {
		secs++
		ps.Add(vs)
		minV, maxV := ps.Extents()
		w := maxV.X - minV.X + 1
		h := maxV.Y - minV.Y + 1

		fmt.Println("-----------")

		if w > 100 || h > 100 {
			fmt.Printf("width: %d,height: %d too big\n", w, h)
			continue
		} else {
			fmt.Println(NewGridWithExtents(ps, minV, maxV, ' ', '#'))
			fmt.Printf("secs: %d  (%d,%d)\n", secs, w, h)
		}

		os.Stdin.Read(b)
	}
}

func getData(filename string) (Vectors, Vectors) {
	lines, _ := file.GetLines(filename)
	ps := make(Vectors, len(lines))
	vs := make(Vectors, len(lines))

	for i, line := range lines {
		matches := reLight.FindStringSubmatch(line)
		ps[i] = Vector{X: atoi(matches[1]), Y: atoi(matches[2])}
		vs[i] = Vector{X: atoi(matches[3]), Y: atoi(matches[4])}
	}

	return ps, vs
}
