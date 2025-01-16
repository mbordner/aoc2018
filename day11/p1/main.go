package main

import "fmt"

type Grid struct {
	sn   int
	data [][]int
}

func (g *Grid) PowerLevel(x, y int) int {
	rackId := x + 10
	pl := rackId * y
	pl += g.sn
	pl *= rackId
	pl = pl % 1000 / 100
	pl -= 5
	return pl
}

func NewGrid(serialNumber int) *Grid {
	grid := &Grid{sn: serialNumber, data: make([][]int, serialNumber)}
	grid.data = make([][]int, 300)
	for j := 0; j < 300; j++ {
		grid.data[j] = make([]int, 300)
		for i := 0; i < 300; i++ {
			grid.data[j][i] = grid.PowerLevel(i+1, j+1)
		}
	}
	return grid
}

func (g *Grid) LargestPowerLevelSubGrid(w, h int) (int, int, int) {
	maxSum := 0
	maxX, maxY := 0, 0

	for y := 0; y < len(g.data)-h; y++ {
		for x := 0; x < len(g.data[y])-w; x++ {
			sum := 0
			for j := y; j < y+h; j++ {
				for i := x; i < x+w; i++ {
					sum += g.data[j][i]
				}
			}
			if sum > maxSum {
				maxSum = sum
				maxX, maxY = x, y
			}
		}
	}

	return maxX + 1, maxY + 1, maxSum
}

func main() {

	g := NewGrid(8199)

	fmt.Println(g.LargestPowerLevelSubGrid(3, 3))
}
