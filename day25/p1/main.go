package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common/file"
	"strconv"
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
					vector := v.Add(Vector{X: x, Y: y, Z: z, W: w})
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
	/*
		vectors := Vector{}.ConstellationNeighborPossibilities()
		fmt.Println(vectors)
		fmt.Println(len(vectors))
	*/
	vectors := getData("../data.txt")
	var constellations []Vectors

	for _, vector := range vectors {
		var in, out []Vectors
		for _, constellation := range constellations {
			inConstellation := false
			for _, v := range constellation {
				if v.Dis(vector) <= 3 {
					inConstellation = true
					break
				}
			}

			if inConstellation {
				in = append(in, constellation)
			} else {
				out = append(out, constellation)
			}
		}

		constellations = out

		if len(in) == 0 {
			in = append(in, Vectors{vector})
		} else {
			tmp := Vectors{vector}
			for _, ic := range in {
				tmp = append(tmp, ic...)
			}
			in = []Vectors{tmp}
		}

		constellations = append(constellations, in...)
	}

	fmt.Println(len(constellations))
}

func getData(filename string) Vectors {
	lines, _ := file.GetLines(filename)
	vectors := make(Vectors, len(lines))
	for i, line := range lines {
		tokens := strings.Split(line, ",")
		vectors[i] = Vector{X: atoi(tokens[0]), Y: atoi(tokens[1]), Z: atoi(tokens[2]), W: atoi(tokens[3])}
	}
	return vectors
}

func atoi(s string) int {
	val, _ := strconv.ParseInt(s, 10, 64)
	return int(val)
}
