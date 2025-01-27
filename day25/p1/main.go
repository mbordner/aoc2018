package main

import (
	"fmt"
	"strings"
)

type Vector struct {
	X, Y, Z, W int
}

func (v Vector) String() string {
	return fmt.Sprintf("(%d,%d,%d,%d)", v.X, v.Y, v.Z, v.W)
}

func (v Vector) Add(o Vector) Vector {
	return Vector{X: v.X + o.X, Y: v.Y + o.Y, Z: v.Z + o.Z, W: v.W}
}

func (v Vector) Dis(o Vector) int {
	return abs(o.X-v.X) + abs(o.Y-v.Y) + abs(o.Z-v.Z) + abs(o.W-v.W)
}

func (v Vector) ConstellationNeighborPossibilities() Vectors {
	var vectors Vectors
	for x := -3; x <= 3; x++ {
		for y := -3 + abs(x); y <= 3-abs(x); y++ {
			for z := -3 + abs(x) + abs(y); z <= 3-abs(x)-abs(y); z++ {
				for w := -3 + abs(x) + abs(y) + abs(z); w <= 3-abs(x)-abs(y)-abs(z); w++ {
					vector := Vector{X: x, Y: y, Z: z, W: w}
					if vector != v {
						vectors = append(vectors, vector)
					}
				}
			}
		}
	}
	return vectors
}

type Vectors []Vector

func (vs Vectors) String() string {
	sb := strings.Builder{}
	sb.WriteString("[\n")
	for _, v := range vs {
		sb.WriteString("  ")
		sb.WriteString(v.String())
		sb.WriteString(",\n")
	}
	sb.WriteString("]\n")
	return sb.String()
}

func abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

func main() {
	vectors := Vector{}.ConstellationNeighborPossibilities()
	fmt.Println(vectors)
	fmt.Println(len(vectors))
}
