package hexagon

import "strings"

type IntNumber interface {
	int | int32 | int64
}

// Cell https://www.redblobgames.com/grids/hexagons/
type Cell[T IntNumber] struct {
	Q T
	R T
	S T
}

func Abs[T IntNumber](x T) T {
	if x >= T(0) {
		return x
	}
	return -x
}

func (c Cell[T]) Distance(o Cell[T]) T {
	return (Abs(c.Q-o.Q) + Abs(c.R-o.R) + Abs(c.S-o.S)) / T(2)
}

type CellLinker[T IntNumber] map[Cell[T]]Cell[T]

func (cl CellLinker[T]) Link(c, o Cell[T]) {
	cl[c] = o
}

func (cl CellLinker[T]) HasLink(c Cell[T]) bool {
	if _, e := cl[c]; e {
		return true
	}
	return false
}

type Grid[T IntNumber] map[Cell[T]]bool

func (g Grid[T]) Has(c Cell[T]) bool {
	if _, e := g[c]; e {
		return true
	}
	return false
}

type Cells[T IntNumber] []Cell[T]

type Dir int

type Dirs []Dir

func DirsFromString(s string) Dirs {
	ds := strings.Split(strings.Join(strings.Fields(s), ""), ",")
	dirs := make(Dirs, len(ds))
	for i, d := range ds {
		dirs[i] = StrDir(d)
	}
	return dirs
}

func StrDir(s string) Dir {
	d := strings.TrimSpace(s)
	d = strings.ToLower(d)
	switch d {
	case "n":
		return N
	case "ne":
		return NE
	case "se":
		return SE
	case "s":
		return S
	case "sw":
		return SW
	case "nw":
		return NW
	}
	return UNKNOWN
}

const (
	UNKNOWN Dir = iota
	N
	NE
	SE
	S
	SW
	NW
)

var (
	Directions = []Dir{N, NE, SE, S, SW, NW}
)

func (c Cell[T]) Valid() bool {
	return c.Q+c.R+c.S == 0
}

func (c Cell[T]) Next(d Dir) Cell[T] {
	switch d {
	case N:
		return Cell[T]{Q: c.Q, R: c.R - 1, S: c.S + 1}
	case NE:
		return Cell[T]{Q: c.Q + 1, R: c.R - 1, S: c.S}
	case SE:
		return Cell[T]{Q: c.Q + 1, R: c.R, S: c.S - 1}
	case S:
		return Cell[T]{Q: c.Q, R: c.R + 1, S: c.S - 1}
	case SW:
		return Cell[T]{Q: c.Q - 1, R: c.R + 1, S: c.S}
	case NW:
		return Cell[T]{Q: c.Q - 1, R: c.R, S: c.S + 1}
	}
	return c
}

func (c Cell[T]) Neighbors() Cells[T] {
	ns := make(Cells[T], len(Directions))
	for i, d := range Directions {
		ns[i] = c.Next(d)
	}
	return ns
}
