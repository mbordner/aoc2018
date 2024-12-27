package geom

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Vector = Point
type Points []Point
type Vectors []Vector
type Cuboids []Cuboid

type Axis int

const (
	X Axis = iota
	Y
	Z
)

type Corner int

const ( //  Left, Right,  Front, Back, Ground, Top
	LFG Corner = iota // Left (x min), Front (y min),  Ground (z min)
	RFG               // Right (x max), Front (y min), Ground (z min)
	LBG               // Left (x min), Back (y max), Ground (z min)
	RBG               // Right (x max), Back (y max), Ground (z min)

	LFT // Left (x min), Front (y min), Top (z max)
	RFT // Right (x max), Front (y min), Top (z max)
	LBT // Left (x min), Back (y max), Top (z max)
	RBT // Right (x max), Back (y max), Top (z max)
)

func (ps Points) Contains(p Point) bool {
	for _, tp := range ps {
		if tp == p {
			return true
		}
	}
	return false
}

func (cs Cuboids) Contains(c Cuboid) bool {
	for _, tc := range cs {
		if tc == c {
			return true
		}
	}
	return false
}

func (c Cuboid) Clone() Cuboid {
	return Cuboid{Min: c.Min.Clone(), Max: c.Max.Clone()}
}

func (cs Cuboids) Clone() Cuboids {
	ns := make(Cuboids, 0, len(cs))
	for _, c := range cs {
		ns = append(ns, c.Clone())
	}
	return ns
}

func (c Cuboid) Points() Points {
	ps := make(Points, 0, c.PointsCount())
	for x := c.Min.X; x <= c.Max.X; x++ {
		for y := c.Min.Y; y <= c.Max.Y; y++ {
			for z := c.Min.Z; z <= c.Max.Z; z++ {
				ps = append(ps, Point{X: x, Y: y, Z: z})
			}
		}
	}
	return ps
}

func (cs Cuboids) Points() Points {
	ps := make(Points, 0, cs.PointsCount())
	for _, c := range cs {
		ps = append(ps, c.Points()...)
	}
	return ps.DeDup().Sort()
}

func (ps Points) Sort() Points {
	sort.Slice(ps, func(i, j int) bool {
		if (ps)[i].X == ps[j].X && ps[i].Y == ps[j].Y {
			if ps[i].Z < ps[j].Z {
				return true
			}
		}
		if ps[i].X == ps[j].X {
			if ps[i].Y < ps[j].Y {
				return true
			}
		}
		if ps[i].X < ps[j].X {
			return true
		}
		return false
	})
	return ps
}

func (ps Points) String() string {
	ss := make([]string, len(ps), len(ps))
	for i, p := range ps {
		ss[i] = p.String()
	}
	return strings.Join(ss, ",")
}

func (ps Points) DeDup() Points {
	tmap := make(map[Point]bool)
	nps := make(Points, 0, len(ps))
	for _, p := range ps {
		if _, exists := tmap[p]; !exists {
			nps = append(nps, p)
			tmap[p] = true
		}
	}
	return nps
}

func (ps Points) Print() {
	for _, p := range ps {
		fmt.Println(p)
	}
}

func (cs Cuboids) Volume() uint64 {
	v := uint64(0)
	for _, c := range cs {
		v += c.Volume()
	}
	return v
}

func (cs Cuboids) PointsCount() uint64 {
	v := uint64(0)
	for _, c := range cs {
		v += c.PointsCount()
	}
	return v
}

func (cs Cuboids) DeDup() Cuboids {
	tmap := make(map[Cuboid]bool)
	ncs := make(Cuboids, 0, len(cs))
	for _, c := range cs {
		if _, exists := tmap[c]; !exists {
			tmap[c] = true
			ncs = append(ncs, c)
		}
	}
	return ncs
}

func (cs Cuboids) Merge(o Cuboid) Cuboids {

	if len(cs) > 0 {
		onlyCS := make(Cuboids, 0, len(cs)+100)

		for _, c := range cs {
			tcs, _, _ := c.Intersect(o)
			if len(tcs) > 0 {
				onlyCS = append(onlyCS, tcs...)
			}
		}

		ns := append(Cuboids{o}, onlyCS.Combine()...)

		return ns.Combine()
	}

	return Cuboids{o}

}

func (cs Cuboids) Remove(o Cuboid) Cuboids {

	if len(cs) > 0 {
		onlyCS := make(Cuboids, 0, len(cs)+100)

		for _, c := range cs {
			tcs, _, _ := c.Intersect(o)
			if len(tcs) > 0 {
				onlyCS = append(onlyCS, tcs...)
			}
		}

		return onlyCS.Combine()
	}

	return Cuboids{}

}

func (cs Cuboids) Combine() Cuboids {
	if len(cs) < 2 {
		return cs
	}

	pmap := make(map[Point]*Cuboid)
	cmap := make(map[*Cuboid]bool)

outer:
	for i := 0; i < len(cs); i++ {
		for j := 0; j < len(cs); j++ {
			if i != j {
				if cs[j].Encloses(cs[i]) {
					continue outer
				}
			}
		}
		cmap[&cs[i]] = true
	}

	psExist := func(ps Points, corners []Corner) *Cuboid {
		var cp *Cuboid = nil
		for _, p := range ps {
			if _, exists := pmap[p]; exists {
				if cp != nil && cp != pmap[p] {
					return nil
				}
				cp = pmap[p]
			} else {
				return nil
			}
		}
		if cp != nil {
			if _, exists := cmap[cp]; exists {
				// we found another cuboid that matches these corners, but in the case where we are only
				// matching 2 distinct points, we need to ensure the matching cuboid only has 4 distinct corners
				checkCorners := cp.Corners()
				for _, crn := range corners {
					if ps.Contains(checkCorners[crn]) == false {
						return nil
					}
				}

				return cp
			}
		}
		return nil
	}

	updatePmap := func() {
		pmap = make(map[Point]*Cuboid)
		for t := range cmap {
			for _, p := range t.Corners() {
				pmap[p] = t
			}
		}
	}

	merged := true
	for merged {
		merged = false

		updatePmap()

		ns := make([]*Cuboid, 0, len(cmap))
		for t := range cmap {
			ns = append(ns, t)
		}
		for _, c := range ns {

			if _, exists := cmap[c]; !exists {
				continue
			}

			corners := c.Corners()
			// left
			if cp := psExist(Points{corners[LFG], corners[LFT], corners[LBG], corners[LFT]}.Transform(Vector{X: -1}), []Corner{RFG, RFT, RBG, RFT}); cp != nil {
				c.Min = (*cp).Min.Clone()
				delete(cmap, cp)
				updatePmap()
				merged = true
				continue
			}

			// right
			if cp := psExist(Points{corners[RFG], corners[RFT], corners[RBG], corners[RFT]}.Transform(Vector{X: 1}), []Corner{LFG, LFT, LBG, LFT}); cp != nil {
				c.Max = (*cp).Max.Clone()
				delete(cmap, cp)
				updatePmap()
				merged = true
				continue
			}

			// front
			if cp := psExist(Points{corners[LFG], corners[RFG], corners[LFT], corners[RFT]}.Transform(Vector{Y: -1}), []Corner{LBG, RBG, LBT, RBT}); cp != nil {
				c.Min = (*cp).Min.Clone()
				delete(cmap, cp)
				updatePmap()
				merged = true
				continue
			}

			// back
			if cp := psExist(Points{corners[LBG], corners[RBG], corners[LBT], corners[RBT]}.Transform(Vector{Y: 1}), []Corner{LFG, RFG, LFT, RFT}); cp != nil {
				c.Max = (*cp).Max.Clone()
				delete(cmap, cp)
				updatePmap()
				merged = true
				continue
			}

			// ground
			if cp := psExist(Points{corners[LBG], corners[RBG], corners[LFG], corners[RFG]}.Transform(Vector{Z: -1}), []Corner{LBT, RBT, LFT, RFT}); cp != nil {
				c.Min = (*cp).Min.Clone()
				delete(cmap, cp)
				updatePmap()
				merged = true
				continue
			}

			// top
			if cp := psExist(Points{corners[LBT], corners[RBT], corners[LFT], corners[RFT]}.Transform(Vector{Z: 1}), []Corner{LBG, RBG, LFG, RFG}); cp != nil {
				c.Max = (*cp).Max.Clone()
				delete(cmap, cp)
				updatePmap()
				merged = true
				continue
			}
		}

	}

	ns := make(Cuboids, 0, len(cmap))
	for c := range cmap {
		ns = append(ns, *c)
	}

	return ns
}

func (cs Cuboids) BreakOverlaps() Cuboids {
	ns := make(Cuboids, 0, len(cs)+100)

	for i := 0; i < len(cs); i++ {
		ns = ns.Merge(cs[i])
	}

	return ns
}

func NewPoints(s string) Points {
	pss := strings.Split(s, ",")
	points := make(Points, 0, len(pss)/3)
	for len(pss) > 0 {
		points = append(points, NewPoint(strings.Join(pss[0:3], ",")))
		pss = pss[3:]
	}
	return points
}

type Point struct {
	X int64
	Y int64
	Z int64
}

func (p Point) String() string {
	return fmt.Sprintf("%d,%d,%d", p.X, p.Y, p.Z)
}

func (p Point) Clone() Point {
	return Point{X: p.X, Y: p.Y, Z: p.Z}
}

func NewPoint(s string) Point {
	p := Point{}
	tokens := strings.Split(s, ",")
	for i := range tokens {
		val, _ := strconv.ParseInt(tokens[i], 10, 64)
		switch i {
		case 0:
			p.X = val
		case 1:
			p.Y = val
		case 2:
			p.Z = val
		}
	}
	return p
}

func NewVector(s string) Vector {
	return Vector(NewPoint(s))
}

func (p Point) Transform(v Vector) Point {
	return Point{X: p.X + v.X, Y: p.Y + v.Y, Z: p.Z + v.Z}
}

func (ps Points) Transform(v Vector) Points {
	ns := make(Points, len(ps), len(ps))
	for i := range ps {
		ns[i] = ps[i].Transform(v)
	}
	return ns
}

func abs(val int64) uint64 {
	if val > 0 {
		return uint64(val)
	}
	return uint64(-val)
}

func min(i, j int64) int64 {
	if i < j {
		return i
	}
	return j
}

func max(i, j int64) int64 {
	if i > j {
		return i
	}
	return j
}

type Cuboid struct {
	Min Point
	Max Point
}

func (c Cuboid) String() string {
	return fmt.Sprintf("%s,%s", c.Min, c.Max)
}

func NewCuboid(s string) Cuboid {
	tokens := strings.Split(s, ",")
	c := Cuboid{}
	c.Min = NewPoint(strings.Join(tokens[0:3], ","))
	c.Max = NewPoint(strings.Join(tokens[3:], ","))
	return c
}

func (c Cuboid) Corners() Points {
	return []Point{
		{c.Min.X, c.Min.Y, c.Min.Z}, // LFG
		{c.Max.X, c.Min.Y, c.Min.Z}, // RFG
		{c.Min.X, c.Max.Y, c.Min.Z}, // LBG
		{c.Max.X, c.Max.Y, c.Min.Z}, // RBG

		{c.Min.X, c.Min.Y, c.Max.Z}, // LFT
		{c.Max.X, c.Min.Y, c.Max.Z}, // RFT
		{c.Min.X, c.Max.Y, c.Max.Z}, // LBT
		{c.Max.X, c.Max.Y, c.Max.Z}, // RBT
	}
}

func (c Cuboid) Contains(p Point) bool {
	if p.X < c.Min.X || p.X > c.Max.X {
		return false
	}
	if p.Y < c.Min.Y || p.Y > c.Max.Y {
		return false
	}
	if p.Z < c.Min.Z || p.Z > c.Max.Z {
		return false
	}

	return true
}

func (c Cuboid) Volume() uint64 {
	xs := abs(c.Max.X - c.Min.X)
	ys := abs(c.Max.Y - c.Min.Y)
	zs := abs(c.Max.Z - c.Min.Z)
	return xs * ys * zs
}

func (c Cuboid) PointsCount() uint64 {
	xs := abs(c.Max.X-c.Min.X) + 1
	ys := abs(c.Max.Y-c.Min.Y) + 1
	zs := abs(c.Max.Z-c.Min.Z) + 1
	return xs * ys * zs
}

func (c Cuboid) IntersectingCorners(o Cuboid) Points {
	pts := make(Points, 0, 8)
	for _, p := range o.Corners() {
		if c.Contains(p) {
			pts = append(pts, p)
		}
	}
	return pts
}

func (c Cuboid) IsCorner(p Point) bool {
	for _, c := range c.Corners() {
		if c == p {
			return true
		}
	}
	return false
}

func (c Cuboid) IsOnEdge(p Point) bool {
	if c.Contains(p) {
		if p.X == c.Min.X || p.X == c.Max.X ||
			p.Y == c.Min.Y || p.Y == c.Max.Y ||
			p.Z == c.Min.Z || p.Z == c.Max.Z {
			return true
		}
	}
	return false
}

func (c Cuboid) Transform(v Vector) Cuboid {
	return Cuboid{Min: c.Min.Transform(v), Max: c.Max.Transform(v)}
}

func (c Cuboid) SplitAt(p Point) Cuboids {
	if !c.Contains(p) {
		return Cuboids{}
	}

	pPP := p.Transform(Vector{1, 1, 1})

	min := Cuboid{c.Min, p}
	max := Cuboid{pPP, c.Max}

	candidates := Cuboids{
		// LFG
		min,

		// RFG
		Cuboid{Min: min.Min.Transform(Vector{X: pPP.X - min.Min.X}), Max: pPP.Transform(Vector{X: max.Max.X - pPP.X, Y: -1, Z: -1})},

		// LBG
		Cuboid{Min: min.Min.Transform(Vector{Y: pPP.Y - min.Min.Y}), Max: pPP.Transform(Vector{X: -1, Y: max.Max.Y - pPP.Y, Z: -1})},

		// RBG
		Cuboid{Min: pPP.Transform(Vector{Z: -(pPP.Z - min.Min.Z)}), Max: pPP.Transform(Vector{X: max.Max.X - pPP.X, Y: max.Max.Y - pPP.Y, Z: -1})},

		// LFT
		Cuboid{Min: min.Min.Transform(Vector{Z: pPP.Z - min.Min.Z}), Max: p.Transform(Vector{Z: max.Max.Z - p.Z})},

		// RFT
		Cuboid{Min: pPP.Transform(Vector{Y: -(pPP.Y - min.Min.Y)}), Max: max.Max.Transform(Vector{Y: -(max.Max.Y - pPP.Y + 1)})},

		// LBT
		Cuboid{Min: pPP.Transform(Vector{X: -(pPP.X - min.Min.X)}), Max: max.Max.Transform(Vector{X: -(max.Max.X - min.Max.X)})},

		// RBT
		max,
	}

	return c.DiscardNonEnclosing(candidates)
}

func (c Cuboid) DiscardNonEnclosing(candidates Cuboids) Cuboids {
	cuboids := make(Cuboids, 0, 8)
	for _, check := range candidates {
		if c.Encloses(check) {
			cuboids = append(cuboids, check)
		}
	}
	return cuboids
}

func (c Cuboid) Encloses(o Cuboid) bool {
	for _, pt := range o.Corners() {
		if !c.Contains(pt) {
			return false
		}
	}
	return true
}

func (cs Cuboids) Overlaps() *Cuboid {
	if len(cs) > 0 {
		c := cs[0]
		for _, o := range cs[1:] {
			c = Cuboid{Min: c.Min.Max(o.Min), Max: c.Max.Min(o.Max)}
		}
		return &c
	}
	return nil
}

func (c Cuboid) Overlaps(o Cuboid) bool {
	np := Cuboids{c, o}.Overlaps()
	if np != nil {
		n := *np
		if c.Encloses(n) {
			return true
		}
		if o.Encloses(n) {
			return true
		}
	}
	return false
}

func (p Point) Max(o Point) Point {
	return Point{X: max(p.X, o.X), Y: max(p.Y, o.Y), Z: max(p.Z, o.Z)}
}

func (p Point) Min(o Point) Point {
	return Point{X: min(p.X, o.X), Y: min(p.Y, o.Y), Z: min(p.Z, o.Z)}
}

func (p Point) Snap(o Point, a Axis) Point {
	n := p.Clone()
	switch a {
	case X:
		n.X = o.X
	case Y:
		n.Y = o.Y
	case Z:
		n.Z = o.Z
	}
	return n
}

func (c Cuboid) Intersect(o Cuboid) (cCuboids Cuboids, intersecting Cuboids, oCuboids Cuboids) {
	if c == o {
		return []Cuboid{}, []Cuboid{c}, []Cuboid{}
	}
	if !c.Overlaps(o) {
		return []Cuboid{c}, []Cuboid{}, []Cuboid{o}
	}

	icptr := Cuboids{c, o}.Overlaps()
	if icptr != nil {

		intersecting = Cuboids{*icptr}

		oCuboids = Cuboids{}

		corners := (*icptr).Corners()

		candidates := make(Cuboids, 0, 8)

		// left side
		candidates = append(candidates, Cuboid{Min: c.Min, Max: corners[LBT].Snap(c.Max, Y).Snap(c.Max, Z).Transform(Vector{X: -1})})
		// right side
		candidates = append(candidates, Cuboid{Min: corners[RFG].Snap(c.Min, Y).Snap(c.Min, Z).Transform(Vector{X: 1}), Max: c.Max})

		// bottom side
		candidates = append(candidates, Cuboid{Min: corners[LFG].Snap(c.Min, Z), Max: corners[RBG].Transform(Vector{Z: -1})})

		// top side
		candidates = append(candidates, Cuboid{Min: corners[LFT].Transform(Vector{Z: 1}), Max: corners[RBT].Snap(c.Max, Z)})

		// front side
		candidates = append(candidates, Cuboid{Min: corners[LFG].Snap(c.Min, Y).Snap(c.Min, Z), Max: corners[RFT].Snap(c.Max, Z).Transform(Vector{Y: -1})})
		// back side
		candidates = append(candidates, Cuboid{Min: corners[LBG].Snap(c.Min, Z).Transform(Vector{Y: 1}), Max: corners[RBT].Snap(c.Max, Y).Snap(c.Max, Z)})

		cCuboids = c.DiscardNonEnclosing(candidates)

		if o != *icptr {
			oCuboids, _, _ = o.Intersect(*icptr)
		}

	}

	return
}
