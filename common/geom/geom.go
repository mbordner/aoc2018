package geom

import (
	"fmt"
	"github.com/mbordner/aoc2018/common/array"
	"log"
	"strings"
)

type Direction int

type IntNumber interface {
	int | int32 | int64
}

type GridBox[T IntNumber] struct {
	P0 Pos[T] // nw corner
	P1 Pos[T] // ne corner
	P2 Pos[T] // se corner
	P3 Pos[T] // sw corner
}

func NewGridBox[T IntNumber](p0, p1, p2, p3 Pos[T]) GridBox[T] {
	gb := GridBox[T]{
		P0: p0,
		P1: p1,
		P2: p2,
		P3: p3,
	}
	if (p0.Y != p1.Y) ||
		(p2.Y != p3.Y) ||
		(p2.Y == p1.Y) ||
		(p0.X != p3.X) ||
		(p1.X != p2.X) ||
		(p1.X == p3.X) {
		log.Fatalln("invalid grid box points")
	}
	return gb
}

func (gb GridBox[T]) Area() T {
	top := GridLine[T]{P0: gb.P0, P1: gb.P1}
	right := GridLine[T]{P0: gb.P1, P1: gb.P2}
	tl := top.Length()
	rl := right.Length()
	return tl * rl
}

type PosGridLines[T IntNumber] map[Pos[T]][]GridLine[T]
type GridLines[T IntNumber] map[GridLine[T]]Direction

type GridLine[T IntNumber] struct {
	P0 Pos[T]
	P1 Pos[T]
}

func (g GridLine[T]) Length() T {
	o := g.P1.Subtract(g.P0)
	return Abs(o.X + o.Y)
}

func (g GridLine[T]) Direction() Direction {
	o := g.P1.Subtract(g.P0)
	o = o.Normalize()
	switch o {
	case Pos[T]{X: T(0), Y: T(-1), Z: T(0)}:
		return North
	case Pos[T]{X: T(1), Y: T(0), Z: T(0)}:
		return East
	case Pos[T]{X: T(0), Y: T(1), Z: T(0)}:
		return South
	case Pos[T]{X: T(-1), Y: T(0), Z: T(0)}:
		return West
	}
	return Unknown
}

func (g GridLine[T]) ContainsGridLine(gl GridLine[T]) bool {
	return g.ContainsPoint(gl.P0) && g.ContainsPoint(gl.P1)
}

func (g GridLine[T]) ContainsPoint(p Pos[T]) bool {
	if p.X == g.P0.X && p.X == g.P1.X {
		if p.Y >= Min(g.P0.Y, g.P1.Y) && p.Y <= Max(g.P0.Y, g.P1.Y) {
			return true
		}
	} else if p.Y == g.P0.Y && p.Y == g.P1.Y {
		if p.X >= Min(g.P0.X, g.P1.X) && p.X <= Max(g.P0.X, g.P1.X) {
			return true
		}
	}
	return false
}

func (pgls PosGridLines[T]) Add(p Pos[T], gl GridLine[T]) {
	if gls, e := pgls[p]; e {
		if !array.Contains(gls, gl) {
			pgls[p] = append(pgls[p], gl)
		}
	} else {
		pgls[p] = []GridLine[T]{gl}
	}
}

func (pgls PosGridLines[T]) AddLine(gl GridLine[T]) {
	pgls.Add(gl.P0, gl)
	pgls.Add(gl.P1, gl)
}

const (
	Unknown Direction = 0
	North   Direction = 1
	South   Direction = 2
	West    Direction = 4
	East    Direction = 8
)

func (d Direction) Is(dirs []Direction) bool {
	for _, dir := range dirs {
		if int(d)&int(dir) == int(dir) {
			return true
		}
	}
	return false
}

func (d Direction) Not(possible []Direction, exclude []Direction) []Direction {
	can := make(map[Direction]bool)
	for _, dir := range possible {
		can[dir] = true
	}
	can[d] = false
	for _, dir := range exclude {
		can[dir] = false
	}
	ds := make([]Direction, 0, len(possible))
	for dir, b := range can {
		if b {
			ds = append(ds, dir)
		}
	}
	return ds
}

func (d Direction) Opposite() Direction {
	o := 0
	if int(d)&int(North) == int(North) {
		o |= int(South)
	}
	if int(d)&int(West) == int(West) {
		o |= int(East)
	}
	if int(d)&int(South) == int(South) {
		o |= int(North)
	}
	if int(d)&int(East) == int(East) {
		o |= int(West)
	}
	return Direction(o)
}

type BoundingBox[T IntNumber] struct {
	MinX T
	MaxX T
	MinY T
	MaxY T
	MinZ T
	MaxZ T
}

func (bb *BoundingBox[T]) SetExtents(x1, y1, z1, x2, y2, z2 T) {
	bb.MinX = x1
	bb.MinY = y1
	bb.MinZ = z1
	bb.MaxX = x2
	bb.MaxY = y2
	bb.MaxZ = z2
}

func (bb BoundingBox[T]) XMin() T {
	return bb.MinX
}

func (bb BoundingBox[T]) XMax() T {
	return bb.MaxX
}

func (bb BoundingBox[T]) YMin() T {
	return bb.MinY
}

func (bb BoundingBox[T]) YMax() T {
	return bb.MaxY
}

func (bb BoundingBox[T]) ZMin() T {
	return bb.MinZ
}

func (bb BoundingBox[T]) ZMax() T {
	return bb.MaxZ
}

func (bb BoundingBox[T]) String() string {
	p1 := Pos[T]{X: bb.MinX, Y: bb.MinY}
	p2 := Pos[T]{X: bb.MaxX, Y: bb.MaxY}
	return fmt.Sprintf("[%s, %s]", p1, p2)
}

func (bb *BoundingBox[T]) Extend(p Pos[T]) {
	if p.X < bb.MinX {
		bb.MinX = p.X
	}
	if p.X > bb.MaxX {
		bb.MaxX = p.X
	}
	if p.Y > bb.MaxY {
		bb.MaxY = p.Y
	}
	if p.Y < bb.MinY {
		bb.MinY = p.Y
	}
	if p.Z < bb.MinZ {
		bb.MinZ = p.Z
	}
	if p.Z > bb.MaxZ {
		bb.MaxZ = p.Z
	}
}

func (bb *BoundingBox[T]) Contains(p Pos[T]) bool {
	if p.X < bb.MinX || p.X > bb.MaxX {
		return false
	}
	if p.Y < bb.MinY || p.Y > bb.MaxY {
		return false
	}
	if p.Z < bb.MinZ || p.Z > bb.MaxZ {
		return false
	}
	return true
}

func (bb *BoundingBox[T]) Surrounds(obb *BoundingBox[T]) bool {
	if obb.MinX < bb.MinX {
		return false
	}
	if obb.MaxX > bb.MaxX {
		return false
	}
	if obb.MinY < bb.MinY {
		return false
	}
	if obb.MaxY > bb.MaxY {
		return false
	}
	if obb.MinZ < bb.MinZ {
		return false
	}
	if obb.MaxZ > bb.MaxZ {
		return false
	}
	return true
}

func (bb *BoundingBox[T]) GetDirection(p Pos[T]) Direction {
	dir := 0
	if p.X < bb.MinX {
		dir |= int(West)
	}
	if p.X > bb.MaxX {
		dir |= int(East)
	}
	if p.Y > bb.MaxY {
		dir |= int(North)
	}
	if p.Y < bb.MinY {
		dir |= int(South)
	}
	return Direction(dir)
}

func (bb *BoundingBox[T]) Intersects(p1, p2 Pos[T]) bool {
	return false
}

func (bb *BoundingBox[T]) GetPrintLines(defaultChar rune, chars []rune, pss Positions[T]) []string {
	lines := make([]string, 0, bb.MaxY-bb.MinY+1)

	charMap := make(map[Pos[T]]rune)
	if chars != nil && pss != nil && len(chars) == len(pss) {
		for i := range chars {
			pss[i].Z = 0
			charMap[pss[i]] = chars[i]
		}
	}
	for i := bb.MaxY; i >= bb.MinY; i-- {
		l := bb.MaxX - bb.MinX + 1
		lineRunes := make([]rune, l, l)
		for j, c := bb.MinX, 0; j <= bb.MaxX; j, c = j+1, c+1 {
			lineRunes[c] = defaultChar
			if char, exists := charMap[Pos[T]{X: j, Y: i, Z: 0}]; exists {
				lineRunes[c] = char
			}
		}
		lines = append(lines, string(lineRunes))
	}

	return lines
}

func (bb *BoundingBox[T]) DistanceFromEdge(p Pos[T]) T {
	d := bb.MaxX - p.X

	t := p.X - bb.MinX
	if t < d {
		d = t
	}

	t = p.Y - bb.MinY
	if t < d {
		d = t
	}

	t = bb.MaxY - p.Y
	if t < d {
		d = t
	}

	return d
}

type Positions[T IntNumber] []Pos[T]

type Pos[T IntNumber] struct {
	X T
	Y T
	Z T
}

func (p Pos[T]) Subtract(o Pos[T]) Pos[T] {
	return p.Transform(-o.X, -o.Y, -o.Z)
}

func (p Pos[T]) Scale(s T) Pos[T] {
	return Pos[T]{X: p.X * s, Y: p.Y * s, Z: p.Z * s}
}

func (p Pos[T]) Normalize() Pos[T] {
	o := Pos[T]{X: p.X, Y: p.Y, Z: p.Z}
	if o.X != 0 {
		o.X = o.X / Abs(o.X)
	}
	if o.Y != 0 {
		o.Y = o.Y / Abs(o.Y)
	}
	if o.Z != 0 {
		o.Z = o.Z / Abs(o.Z)
	}
	return o
}

func (p Pos[T]) Transform(x, y, z T) Pos[T] {
	return Pos[T]{X: p.X + x, Y: p.Y + y, Z: p.Z + z}
}

func (p Pos[T]) TransformDir(d Direction, count T) Pos[T] {
	switch d {
	case North:
		return p.Transform(0, T(-count), 0)
	case East:
		return p.Transform(T(count), 0, 0)
	case South:
		return p.Transform(0, T(count), 0)
	case West:
		return p.Transform(T(-count), 0, 0)
	}
	return p.Transform(0, 0, 0)
}

func (p Pos[T]) TransformDirs(d Direction) []Pos[T] {
	np := make([]Pos[T], 0, 4)
	if d.Is([]Direction{North}) {
		np = append(np, p.Transform(0, -1, 0))
	}
	if d.Is([]Direction{East}) {
		np = append(np, p.Transform(1, 0, 0))
	}
	if d.Is([]Direction{South}) {
		np = append(np, p.Transform(0, 1, 0))
	}
	if d.Is([]Direction{West}) {
		np = append(np, p.Transform(-1, 0, 0))
	}
	return np
}
func (p Pos[T]) Diff(o Pos[T]) Pos[T] {
	v := Pos[T]{
		X: p.X - o.X,
		Y: p.Y - o.Y,
		Z: p.Z - o.Z,
	}
	return v
}

func (p Pos[T]) ManhattanDistance(o Pos[T]) T {
	return Max[T](p.X, o.X) - Min[T](p.X, o.X) +
		Max[T](p.Y, o.Y) - Min[T](p.Y, o.Y) +
		Max[T](p.Z, o.Z) - Min[T](p.Z, o.Z)
}

func (p Pos[T]) GetXYPositionsAtManhattanDistance(d T) Positions[T] {
	ps := make(Positions[T], 0, ((d-1)*4)+4)

	left := Pos[T]{X: p.X - d, Y: p.Y}
	right := Pos[T]{X: p.X + d, Y: p.Y}
	top := Pos[T]{X: p.X, Y: p.Y - d}
	bottom := Pos[T]{X: p.X, Y: p.Y + d}

	ps = append(ps, Positions[T]{top, bottom, left, right}...)

	for j := T(1); j < d; j++ {
		dx := d - j
		topY, bottomY := p.Y-j, p.Y+j
		topLeft := Pos[T]{X: p.X - dx, Y: topY}
		topRight := Pos[T]{X: p.X + dx, Y: topY}
		bottomLeft := Pos[T]{X: p.X - dx, Y: bottomY}
		bottomRight := Pos[T]{X: p.X + dx, Y: bottomY}

		ps = append(ps, Positions[T]{topLeft, topRight, bottomLeft, bottomRight}...)
	}

	return ps
}

// GetXYPositionsWithinManhattanDistance returns positions within distance in x,y plane
func (p Pos[T]) GetXYPositionsWithinManhattanDistance(d T) Positions[T] {
	pm := make(map[Pos[T]]bool)

	for j := T(0); j <= d; j++ {
		for i := T(0); i <= d-j; i++ {
			topY, bottomY := p.Y-j, p.Y+j
			topLeft := Pos[T]{X: p.X - i, Y: topY}
			topRight := Pos[T]{X: p.X + i, Y: topY}
			bottomLeft := Pos[T]{X: p.X - i, Y: bottomY}
			bottomRight := Pos[T]{X: p.X + i, Y: bottomY}
			pm[topLeft] = true
			pm[topRight] = true
			pm[bottomLeft] = true
			pm[bottomRight] = true
		}
	}

	ps := make(Positions[T], 0, len(pm))
	for p := range pm {
		ps = append(ps, p)
	}
	return ps
}

func (p Pos[T]) Clone() Pos[T] {
	return Pos[T]{X: p.X, Y: p.Y, Z: p.Z}
}

func (ps Positions[T]) String() string {
	strs := make([]string, 0, len(ps))
	for _, p := range ps {
		strs = append(strs, p.String())
	}
	return strings.Join(strs, ",")
}

func (ps Positions[T]) Contains(p Pos[T]) bool {
	for _, tp := range ps {
		if tp == p {
			return true
		}
	}
	return false
}

func (p Pos[T]) String() string {
	return fmt.Sprintf("{x:%d, y:%d, z:%d}", p.X, p.Y, p.Z)
}

func Min[T IntNumber](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T IntNumber](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Abs[T IntNumber](n T) T {
	if n < 0 {
		return -n
	}
	return n
}

func (bb *BoundingBox[T]) GetPositionsSize() T {
	xs := Abs(bb.MaxX-bb.MinX) + 1
	ys := Abs(bb.MaxY-bb.MinY) + 1
	zs := Abs(bb.MaxZ-bb.MinZ) + 1
	return xs * ys * zs
}

func (bb *BoundingBox[T]) GetPositions() Positions[T] {
	poss := make(Positions[T], 0, ((bb.MaxX-bb.MinX)+1)*((bb.MaxY-bb.MinY)+1*((bb.MaxZ-bb.MinZ)+1)))
	for z := bb.MinZ; z <= bb.MaxZ; z++ {
		for y := bb.MinY; y <= bb.MaxY; y++ {
			for x := bb.MinX; x <= bb.MaxX; x++ {
				poss = append(poss, Pos[T]{Z: z, Y: y, X: x})
			}
		}
	}
	return poss
}

func (ps *Positions[T]) Transform(x, y, z T) Positions[T] {
	for i := 0; i < len(*ps); i++ {
		(*ps)[i].X += x
		(*ps)[i].Y += y
		(*ps)[i].Z += z
	}
	return *ps
}
