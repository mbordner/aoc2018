package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common"
	"github.com/mbordner/aoc2018/common/file"
	"regexp"
	"strconv"
)

var (
	reCave = regexp.MustCompile(`depth:\s+(\d+)\s+target:\s+(\d+),(\d+)`)
)

type RegionType int

const (
	Rocky RegionType = iota
	Wet
	Narrow
)

const (
	RockyChar  = '.'
	WetChar    = '='
	NarrowChar = '|'
)

type Region struct {
	pos           common.Pos
	geologicIndex int
	erosionLevel  int
}

func (r *Region) Type() RegionType {
	return RegionType(r.erosionLevel % 3)
}

type Cave struct {
	depth  int
	target common.Pos
	start  common.Pos
	region map[common.Pos]*Region
}

func (c *Cave) RiskLevel() int {
	rl := 0
	for _, r := range c.region {
		rl += int(r.Type())
	}
	return rl
}

func main() {
	cave := getData("../data.txt")
	fmt.Println(cave.RiskLevel())
}

func getData(filename string) *Cave {
	cave := &Cave{start: common.Pos{}, region: make(map[common.Pos]*Region)}

	content, _ := file.GetContent(filename)
	matches := reCave.FindStringSubmatch(string(content))
	cave.depth = atoi(matches[1])
	cave.target = common.Pos{X: atoi(matches[2]), Y: atoi(matches[3])}

	for y := 0; y <= cave.target.Y; y++ {
		for x := 0; x <= cave.target.X; x++ {
			p := common.Pos{X: x, Y: y}
			r := &Region{pos: p}

			if p == cave.start || p == cave.target {
				r.geologicIndex = 0
			} else if y == 0 {
				r.geologicIndex = x * 16807
			} else if x == 0 {
				r.geologicIndex = y * 48271
			} else {
				r.geologicIndex = cave.region[common.Pos{X: x - 1, Y: y}].erosionLevel *
					cave.region[common.Pos{X: x, Y: y - 1}].erosionLevel
			}

			r.erosionLevel = (r.geologicIndex + cave.depth) % 20183

			cave.region[p] = r
		}
	}

	return cave
}

func atoi(s string) int {
	val, _ := strconv.ParseInt(s, 10, 64)
	return int(val)
}
