package main

import (
	"cmp"
	"fmt"
	"github.com/mbordner/aoc2018/common"
	"github.com/mbordner/aoc2018/common/datastructure"
	"github.com/mbordner/aoc2018/common/file"
	"math"
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

type EquippedGear int

const (
	Neither EquippedGear = iota
	Torch
	ClimbingGear
)

func (eg EquippedGear) Is(g EquippedGear) bool {
	return eg == g
}

const (
	RockyChar  = '.'
	WetChar    = '='
	NarrowChar = '|'
	StartChar  = 'M'
	TargetChar = 'T'
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

func (c *Cave) String() string {
	return c.Grid().String()
}

func (c *Cave) Grid() common.Grid {
	maxX, maxY := c.start.X, c.start.Y
	for p := range c.region {
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}
	w := maxX + 1
	h := maxY + 1
	grid := make(common.Grid, h)
	for y := range grid {
		grid[y] = make([]byte, w)
		for x := range grid[y] {
			t := c.Region(x, y).Type()
			switch t {
			case Rocky:
				grid[y][x] = RockyChar
			case Wet:
				grid[y][x] = WetChar
			case Narrow:
				grid[y][x] = NarrowChar
			}
		}
	}
	grid[c.start.Y][c.start.X] = StartChar
	grid[c.target.Y][c.target.X] = TargetChar
	return grid
}

func (c *Cave) PosRegion(pos common.Pos) *Region {
	return c.Region(pos.X, pos.Y)
}

func (c *Cave) Region(x, y int) *Region {
	if x >= 0 && y >= 0 {
		if r, e := c.region[common.Pos{X: x, Y: y}]; e {
			return r
		}
		r := &Region{pos: common.Pos{X: x, Y: y}}

		if y == 0 {
			r.geologicIndex = x * 16807
		} else if x == 0 {
			r.geologicIndex = y * 48271
		} else {
			r.geologicIndex = c.Region(x-1, y).erosionLevel * c.Region(x, y-1).erosionLevel
		}

		r.erosionLevel = (r.geologicIndex + c.depth) % 20183
		c.region[r.pos] = r

		return r
	}
	return nil
}

type Node struct {
	gear EquippedGear
	pos  common.Pos
}

type NodeWithDistance struct {
	node Node
	dis  int
}

func (cur Node) Adjacent(cave *Cave) ([]Node, []int) {
	nodes := make([]Node, 0, 16)
	distances := make([]int, 0, 16)
	for _, d := range common.AdjacentDirs {
		np := cur.pos.Add(d)
		nr := cave.PosRegion(np)
		if nr != nil {
			nt := nr.Type()
			ct := cave.PosRegion(cur.pos).Type()

			if ct == Rocky { // must use Torch or Climbing, not Neither
				if nt == Wet { // can use Neither or Climbing, not Torch
					if cur.gear.Is(Torch) {
						// must change to Climbing only
						nodes = append(nodes, Node{gear: ClimbingGear, pos: np})
						distances = append(distances, 8)
					} else if cur.gear.Is(ClimbingGear) {
						// must not make any changes
						nodes = append(nodes, Node{gear: ClimbingGear, pos: np})
						distances = append(distances, 1)
					} else {
						panic("wrong gear")
					}
				} else if nt == Narrow { // can use Torch or Neither, not Climbing
					if cur.gear.Is(ClimbingGear) {
						// must change to Torch only
						nodes = append(nodes, Node{gear: Torch, pos: np})
						distances = append(distances, 8)
					} else if cur.gear.Is(Torch) {
						// must not make any changes
						nodes = append(nodes, Node{gear: Torch, pos: np})
						distances = append(distances, 1)
					} else {
						panic("wrong gear")
					}
				} else { // must use Torch or Climbing, not Neither
					if cur.gear.Is(ClimbingGear) {
						nodes = append(nodes, Node{gear: cur.gear, pos: np})
						distances = append(distances, 1)

						nodes = append(nodes, Node{gear: Torch, pos: np})
						distances = append(distances, 8)
					} else {
						nodes = append(nodes, Node{gear: cur.gear, pos: np})
						distances = append(distances, 1)

						nodes = append(nodes, Node{gear: ClimbingGear, pos: np})
						distances = append(distances, 8)
					}
				}
			} else if ct == Wet { // can use Neither or Climbing, not Torch
				if nt == Rocky { // must use Torch or Climbing, not Neither
					if cur.gear.Is(Neither) {
						// must equip Climbing only
						nodes = append(nodes, Node{gear: ClimbingGear, pos: np})
						distances = append(distances, 8)
					} else if cur.gear.Is(ClimbingGear) {
						// must not make any changes
						nodes = append(nodes, Node{gear: ClimbingGear, pos: np})
						distances = append(distances, 1)
					} else {
						panic("wrong gear")
					}
				} else if nt == Narrow { // can use Torch or Neither, not Climbing
					if cur.gear.Is(ClimbingGear) {
						// must switch to Neither
						nodes = append(nodes, Node{gear: Neither, pos: np})
						distances = append(distances, 8)
					} else if cur.gear.Is(Neither) {
						// must not make any changes
						nodes = append(nodes, Node{gear: Neither, pos: np})
						distances = append(distances, 1)
					} else {
						panic("wrong gear")
					}
				} else { // can use Neither or Climbing, not Torch
					if cur.gear.Is(ClimbingGear) {
						nodes = append(nodes, Node{gear: cur.gear, pos: np})
						distances = append(distances, 1)

						nodes = append(nodes, Node{gear: Neither, pos: np})
						distances = append(distances, 8)
					} else {
						nodes = append(nodes, Node{gear: cur.gear, pos: np})
						distances = append(distances, 1)

						nodes = append(nodes, Node{gear: ClimbingGear, pos: np})
						distances = append(distances, 8)
					}
				}
			} else if ct == Narrow { // can use Torch or Neither, not Climbing
				if nt == Rocky { // must use Torch or Climbing, not Neither
					if cur.gear.Is(Neither) {
						// must switch to torch
						nodes = append(nodes, Node{gear: Torch, pos: np})
						distances = append(distances, 8)
					} else if cur.gear.Is(Torch) {
						// must not make any changes
						nodes = append(nodes, Node{gear: Torch, pos: np})
						distances = append(distances, 1)
					} else {
						panic("wrong gear")
					}
				} else if nt == Wet { // can use Neither or Climbing, not Torch
					if cur.gear.Is(Torch) {
						// must switch to Neither
						nodes = append(nodes, Node{gear: Neither, pos: np})
						distances = append(distances, 8)
					} else if cur.gear.Is(Neither) {
						// must not make any changes
						nodes = append(nodes, Node{gear: Neither, pos: np})
						distances = append(distances, 1)
					} else {
						panic("wrong gear")
					}
				} else { // can use Torch or Neither, not Climbing
					if cur.gear.Is(Torch) {
						nodes = append(nodes, Node{gear: cur.gear, pos: np})
						distances = append(distances, 1)

						nodes = append(nodes, Node{gear: Neither, pos: np})
						distances = append(distances, 8)
					} else {
						nodes = append(nodes, Node{gear: cur.gear, pos: np})
						distances = append(distances, 1)

						nodes = append(nodes, Node{gear: Torch, pos: np})
						distances = append(distances, 8)
					}
				}
			}

		}
	}
	return nodes, distances
}

type KnownDistances struct {
	dis  map[Node]int
	prev map[Node]Node
}

func NewKnownDistances() *KnownDistances {
	return &KnownDistances{dis: make(map[Node]int), prev: make(map[Node]Node)}
}

func (kd *KnownDistances) Len() int {
	return len(kd.dis)
}

func (kd *KnownDistances) Add(n Node, dis int, prev Node) {
	if d, e := kd.dis[n]; e {
		if dis <= d {
			kd.dis[n] = dis
			kd.prev[n] = prev
		}
	} else {
		kd.dis[n] = dis
		kd.prev[n] = prev
	}
}

func (kd *KnownDistances) Dis(n Node) int {
	if d, e := kd.dis[n]; e {
		return d
	}
	return math.MaxUint32
}

type Visited map[Node]bool

func (v Visited) Seen(n Node) bool {
	if b, e := v[n]; e {
		return b
	}
	return false
}

func main() {
	cave := getData("../data.txt")

	kd := NewKnownDistances()
	visited := make(Visited)
	startNode := Node{gear: Torch, pos: cave.start}
	targetNode := Node{gear: Torch, pos: cave.target}
	kd.Add(startNode, 0, startNode)

	anyHeap := datastructure.NewAnyHeap[NodeWithDistance](func(a, b NodeWithDistance) int {
		return cmp.Compare(a.dis, b.dis)
	})

	anyHeap.Unshift(NodeWithDistance{node: startNode, dis: 0})

	minDistanceToTarget := math.MaxUint32
	var minTargetNode Node

	for anyHeap.Len() > 0 {
		curNWD := anyHeap.Shift()
		cur := curNWD.node

		if !visited.Seen(cur) {
			visited[cur] = true

			neighbors, distances := cur.Adjacent(cave)
			for i, neighbor := range neighbors {
				neighborDistance := curNWD.dis + distances[i]
				if neighbor == targetNode {
					if neighborDistance < minDistanceToTarget {
						minDistanceToTarget = neighborDistance
						minTargetNode = neighbor
					}
				}
				kd.Add(neighbor, neighborDistance, cur)
				if neighborDistance < minDistanceToTarget {
					anyHeap.Unshift(NodeWithDistance{node: neighbor, dis: neighborDistance})
				}
			}
		}
	}

	//fmt.Println(cave)

	fmt.Println(kd.Len())
	fmt.Println(kd.Dis(minTargetNode))
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
