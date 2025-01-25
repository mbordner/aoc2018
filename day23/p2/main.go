package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common"
	"github.com/mbordner/aoc2018/common/file"
	"regexp"
	"strconv"
)

var (
	reNanoBot = regexp.MustCompile(`pos=<(-?\d+),(-?\d+),(-?\d+)>,\s+r=(\d+)`)
)

type Number interface {
	int | int32 | int64 | float32 | float64
}

type Vector[T Number] struct {
	X, Y, Z T
}

func (v Vector[T]) String() string {
	return fmt.Sprintf("{%v, %v, %v}", v.X, v.Y, v.Z)
}

func (v Vector[T]) Add(o Vector[T]) Vector[T] {
	return Vector[T]{X: v.X + o.X, Y: v.Y + o.Y, Z: v.Z + o.Z}
}

func (v Vector[T]) Scale(s T) Vector[T] {
	return Vector[T]{X: v.X * s, Y: v.Y * s, Z: v.Z * s}
}

func (v Vector[T]) Dis(o Vector[T]) T {
	return abs(o.X-v.X) + abs(o.Y-v.Y) + abs(o.Z-v.Z)
}

type NanoBots[T Number] []*NanoBot[T]
type NanoBot[T Number] struct {
	pos    Vector[T]
	radius T
}

func (n *NanoBot[T]) Contains(v Vector[T]) bool {
	if n.pos.Dis(v) <= n.radius {
		return true
	}
	return false
}

func (n *NanoBot[T]) Points() []Vector[T] {
	pts := make([]Vector[T], 0, int(n.radius*n.radius*n.radius))
	for x := -n.radius; x <= n.radius; x++ {
		for y := -n.radius + abs(x); y <= n.radius-abs(x); y++ {
			for z := -n.radius + abs(x) + abs(y); z <= n.radius-abs(x)-abs(y); z++ {
				pts = append(pts, n.pos.Add(Vector[T]{X: x, Y: y, Z: z}))
			}
		}
	}

	return pts
}

func (n *NanoBot[T]) SurfacePoints() []Vector[T] {
	var pts []Vector[T]
	for x := -n.radius; x <= n.radius; x++ {
		for y := -n.radius + abs(x); y <= n.radius-abs(x); y++ {
			for z := -n.radius + abs(x) + abs(y); z <= n.radius-abs(x)-abs(y); z++ {
				p := n.pos.Add(Vector[T]{X: x, Y: y, Z: z})
				if n.pos.Dis(p) == n.radius {
					pts = append(pts, p)
				}
			}
		}
	}
	return pts
}

func (n *NanoBot[T]) SurfacePointsScaled(scale T) []Vector[T] {
	var pts []Vector[T]
	r := n.radius / scale
	np := Vector[T]{X: n.pos.X / scale, Y: n.pos.Y / scale, Z: n.pos.Z / scale}
	for x := -r; x <= r; x++ {
		for y := -r + abs(x); y <= r-abs(x); y++ {
			for z := -r + abs(x) + abs(y); z <= r-abs(x)-abs(y); z++ {
				p := np.Add(Vector[T]{X: x, Y: y, Z: z})
				if np.Dis(p) == r {
					pts = append(pts, p)
				}
			}
		}
	}
	return pts
}

func (n *NanoBot[T]) CountSurfacePointsScaled(pc *PointCounter[T], scale T) {
	r := n.radius / scale
	np := Vector[T]{X: n.pos.X / scale, Y: n.pos.Y / scale, Z: n.pos.Z / scale}
	for x := -r; x <= r; x++ {
		for y := -r + abs(x); y <= r-abs(x); y++ {
			for z := -r + abs(x) + abs(y); z <= r-abs(x)-abs(y); z++ {
				p := np.Add(Vector[T]{X: x, Y: y, Z: z})
				if np.Dis(p) == r {
					pc.AddPoint(p)
				}
			}
		}
	}
}

func (n *NanoBot[T]) CountSurfacePoints(pc *PointCounter[T]) {
	for x := -n.radius; x <= n.radius; x++ {
		for y := -n.radius + abs(x); y <= n.radius-abs(x); y++ {
			for z := -n.radius + abs(x) + abs(y); z <= n.radius-abs(x)-abs(y); z++ {
				p := n.pos.Add(Vector[T]{X: x, Y: y, Z: z})
				if n.pos.Dis(p) == n.radius {
					pc.AddPoint(p)
				}
			}
		}
	}
}

func (n *NanoBot[T]) CountPoints(pc *PointCounter[T]) {
	for x := -n.radius; x <= n.radius; x++ {
		for y := -n.radius + abs(x); y <= n.radius-abs(x); y++ {
			for z := -n.radius + abs(x) + abs(y); z <= n.radius-abs(x)-abs(y); z++ {
				pc.AddPoint(n.pos.Add(Vector[T]{X: x, Y: y, Z: z}))
			}
		}
	}
}

type PointCounter[T Number] struct {
	points map[Vector[T]]int
	maxC   int
	maxP   Vector[T]
}

func NewPointCounter[T Number]() *PointCounter[T] {
	return &PointCounter[T]{points: make(map[Vector[T]]int)}
}

func (pc *PointCounter[T]) Has(p Vector[T]) bool {
	if c, e := pc.points[p]; e {
		if c >= 0 {
			return true
		}
	}
	return false
}

func (pc *PointCounter[T]) AddPoint(p Vector[T]) {
	count := 1
	if c, e := pc.points[p]; e {
		count = c + 1
	}
	pc.points[p] = count
	if count > pc.maxC {
		pc.maxC = count
		pc.maxP = p
		fmt.Println(pc.String())
	} else if count == pc.maxC {
		o := Vector[T]{}
		if pc.maxP.Dis(o) > p.Dis(o) {
			pc.maxP = p
			fmt.Println(pc.String())
		}
	}
}

func (pc *PointCounter[T]) String() string {
	return fmt.Sprintf("%d %s %d", pc.maxC, pc.maxP, pc.maxP.Dis(Vector[T]{}))
}

func (pc *PointCounter[T]) AddPoints(ps []Vector[T]) {
	for _, p := range ps {
		pc.AddPoint(p)
	}
}

//100000000 <-- too high

// 94095928 not right
// 94095822  too low
// 93913372
// 93999712  too low
// 93999709

// {67, 116, 52} 445 at scale: 400000
// {53, 93, 42} 377 at scale:  500000
// {27, 46, 21} 427 at scale: 1000000

// 738 {26723560, 46400000, 20800000} 93923560  < X
// 809 {26800000, 46323560, 20800000} 93923560  < Y
// 809 {26800000, 46400000, 20876440} 94076440
// 828 {26723560, 46323560, 20866252} 93913372
// 828 {26723521, 46323521, 20866252}
// {26723514, 46323514, 20866252}
// 839 {26298448, 46323514, 20866252} 93488214
// 839 {26298448, 46323514, 20866252} 93488214
// 906 {26298448, 46323514, 21291318} 93913280
// 910 {26723514, 46323514, 21291318} 94338346
// 910 {26723439, 46323439, 21291318} 94338196
// 910 {26723383, 46323383, 21291318} 94338084
// 910 {26723323, 46323323, 21291318} 94337964
// 910 {26722803, 46322803, 21291318} 94336924
// 910 {26718003, 46318003, 21291318} 94327324
// 910 {26663803, 46263803, 21291318} 94218924
// 910 {26563803, 46163803, 21291318} 94018924
// 910 {26658803, 46258803, 21291318} 94208924
// 910 {26658749, 46258749, 21291318} 94208816
// 910 {26658702, 46258702, 21291318} 94208722
// 910 {26657922, 46257922, 21291318} 94207162
// 910 {26653622, 46253622, 21291318} 94198562
// 910 {26657252, 46257252, 21291318} 94205822

// 908 {26781123, 46600000, 21100000} 94481123
// 910 {26781123, 46572441, 21100000} 94453564
// 911 {26659123, 46471441, 21079000} 94209564
// 911 {26615123, 46427441, 21079000} 94121564
// 934 {26604863, 46417401, 21078790} 94101054
// 938 {26604461, 46416994, 21078785} 94100240
// 938 {26604390, 46416923, 21078785} 94100098
// 938 {26604326, 46416859, 21078785} 94099970
// 938 {26603368, 46415901, 21078785} 94098054
// 938 {26602665, 46415198, 21078785} 94096648
// 938 {26602305, 46414838, 21078785} 94095928
// 938 {26602305, 46414838, 21078785} 94095928

// 954 {26794863, 46607401, 21078790} 94481054

// 970 {26794893, 46607431, 21078790} 94481114
// 977 {26794906, 46607439, 21078785} 94481130
func main() {
	nanobots, _ := getData[int]("../data.txt")

	scale := int(100000)
	pc := NewPointCounter[int]()

	/*
		for _, n := range nanobots {
			if abs(n.pos.X) < scale || abs(n.pos.Y) < scale || abs(n.pos.Z) < scale || abs(n.radius) < scale {
				fmt.Println("nanobot ", n)
			}
			n.CountSurfacePointsScaled(pc, scale)
		}
	*/

	pc = NewPointCounter[int]()
	//pc.maxP = Vector[int]{X: 67, Y: 116, Z: 52}
	pc.maxP = Vector[int]{X: 268, Y: 466, Z: 211}

	botsMap := make(map[*NanoBot[int]]bool)

	searchPoint := pc.maxP.Scale(scale)
	//searchPoint.X = 26602252
	//searchPoint.Y = 46202252
	//searchPoint.Z = 21291318
	searchPoint.X = 26794906
	searchPoint.Y = 46607439
	searchPoint.Z = 21078785
	searchPoints := []Vector[int]{searchPoint}
	searchPoints = append(searchPoints, searchPoint.Add(Vector[int]{X: scale, Y: scale, Z: scale}))

	for _, sp := range searchPoints {
		for _, n := range nanobots {
			if n.Contains(sp) {
				botsMap[n] = true
			}
		}
	}

	pc = NewPointCounter[int]()
	queue := make(common.Queue[Vector[int]], 0, 100)
	queue.Enqueue(searchPoint)

	for b := range botsMap {
		if b.Contains(searchPoint) {
			pc.AddPoint(searchPoint)
		}
	}

	for !queue.Empty() {
		cur := *(queue.Dequeue())

		for _, d := range []Vector[int]{{-1, 0, 0}, {1, 0, 0}, {0, 1, 0}, {0, -1, 0}, {0, 0, 1}, {0, 0, -1}} {
			nv := cur.Add(d.Scale(1))
			dFrSP := searchPoint.Dis(nv)
			if !pc.Has(nv) && dFrSP <= scale*10 {
				pc.points[nv] = 0
				for _, b := range nanobots {
					if b.Contains(nv) {
						pc.AddPoint(nv)
					}
				}
				queue.Enqueue(nv)
			}
		}
	}

	fmt.Println(pc.maxC, pc.maxP, pc.maxP.Dis(Vector[int]{}))

}

func getData[T Number](filename string) (NanoBots[T], *NanoBot[T]) {
	lines, _ := file.GetLines(filename)
	bots := make(NanoBots[T], len(lines))
	maxSignalRadius := T(0)
	var strongest *NanoBot[T]
	for i, line := range lines {
		matches := reNanoBot.FindStringSubmatch(line)
		bots[i] = &NanoBot[T]{pos: Vector[T]{
			X: atoi[T](matches[1]),
			Y: atoi[T](matches[2]),
			Z: atoi[T](matches[3]),
		}, radius: atoi[T](matches[4])}

		if bots[i].radius > maxSignalRadius {
			maxSignalRadius = bots[i].radius
			strongest = bots[i]
		}
	}
	return bots, strongest
}

func atoi[T Number](s string) T {
	val, _ := strconv.ParseInt(s, 10, 64)
	return T(val)
}

func abs[T Number](x T) T {
	if x >= 0 {
		return x
	}
	return -x
}
