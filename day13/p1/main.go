package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common/file"
	"sort"
)

var (
	N = Pos{-1, 0}
	E = Pos{0, 1}
	S = Pos{1, 0}
	W = Pos{0, -1}

	turnOptions = []int{-1, 0, 1}
)

const (
	Vert  = '|'
	Hor   = '-'
	CNW   = '/'
	CNE   = '\\'
	Cross = '+'
	Down  = 'v'
	Up    = '^'
	Right = '>'
	Left  = '<'
	Space = ' '
	Coll  = 'X'
)

type Pos struct {
	Y int
	X int
}

func (p Pos) String() string {
	return fmt.Sprintf("{%d,%d}", p.X, p.Y)
}

func (p Pos) Add(o Pos) Pos {
	return Pos{p.Y + o.Y, p.X + o.X}
}

type Cart struct {
	id int
	p  Pos
	d  Pos
	t  int
}

func (c *Cart) Char() byte {
	switch c.d {
	case N:
		return Up
	case E:
		return Right
	case S:
		return Down
	case W:
		return Left
	}
	return Space
}

func (c *Cart) Turn() {
	t := c.t
	c.t++
	if c.t == len(turnOptions) {
		c.t = 0
	}
	switch turnOptions[t] {
	case -1:
		switch c.d {
		case N:
			c.d = W
		case E:
			c.d = N
		case S:
			c.d = E
		case W:
			c.d = S
		}
	case 1:
		switch c.d {
		case N:
			c.d = E
		case E:
			c.d = S
		case S:
			c.d = W
		case W:
			c.d = N
		}
	}
}

type Track [][]byte

func (t Track) Clone() Track {
	ot := make(Track, len(t))
	for y := range t {
		ot[y] = make([]byte, len(t[y]))
		copy(ot[y], t[y])
	}
	return ot
}

func (t Track) Print(carts Carts) {
	ot := t.Clone()
	for _, c := range carts {
		tc := ot[c.p.Y][c.p.X]
		if tc == Up || tc == Down || tc == Right || tc == Left || tc == Coll {
			ot[c.p.Y][c.p.X] = Coll
		} else {
			ot[c.p.Y][c.p.X] = c.Char()
		}
	}
	for y := range ot {
		fmt.Println(string(ot[y]))
	}
}

type Carts []*Cart

type Collision Carts

func (cs Carts) Advance(track Track) []Collision {
	collisionMap := make(map[Pos]Collision)
	var collisions []Collision
	for _, c := range cs {
		if _, e := collisionMap[c.p]; e {
			collisionMap[c.p] = append(collisionMap[c.p], c)
			continue
		}

		np := c.p.Add(c.d)
		c.p = np
		if _, e := collisionMap[np]; e {
			collisionMap[np] = append(collisionMap[np], c)
		} else {
			collisionMap[np] = Collision{c}
		}

		if track[np.Y][np.X] == Cross {
			c.Turn()
		} else if track[np.Y][np.X] == CNE {
			switch c.d {
			case N:
				c.d = W
			case E:
				c.d = S
			case S:
				c.d = E
			case W:
				c.d = N
			}
		} else if track[np.Y][np.X] == CNW {
			switch c.d {
			case N:
				c.d = E
			case E:
				c.d = N
			case S:
				c.d = W
			case W:
				c.d = S
			}
		}
	}

	sort.Slice(cs, func(i, j int) bool {
		if cs[i].p.Y < cs[j].p.Y {
			return true
		} else if cs[i].p.Y == cs[j].p.Y {
			return cs[i].p.X < cs[j].p.X
		}
		return false
	})

	for _, cols := range collisionMap {
		if len(cols) > 1 {
			collisions = append(collisions, cols)
		}
	}

	return collisions
}

func main() {
	track, carts := getData("../data.txt")

	moves := 0
	var collisions []Collision
	for {
		collisions = carts.Advance(track)
		moves++
		track.Print(carts)
		if len(collisions) > 0 {
			break
		}
	}

	fmt.Printf("collided after %d moves\n", moves)
	for _, c := range collisions {
		fmt.Printf("collided at %s\n", c[0].p)
	}

}

func getData(filename string) (Track, Carts) {
	lines, _ := file.GetLines(filename)
	w := 0
	for _, line := range lines {
		if len(line) > w {
			w = len(line)
		}
	}

	track := make(Track, len(lines))
	for y, line := range lines {
		track[y] = make([]byte, w)
		for x := range track[y] {
			track[y][x] = Space
		}
		copy(track[y], []byte(line))
	}

	carts := make(Carts, 0, 10)

	id := 0

	for y := range track {
		for x := range track[y] {
			cart := &Cart{p: Pos{Y: y, X: x}}
			switch track[y][x] {
			case Right:
				track[y][x] = Hor
				cart.d = E
			case Left:
				track[y][x] = Hor
				cart.d = W
			case Up:
				track[y][x] = Vert
				cart.d = N
			case Down:
				track[y][x] = Vert
				cart.d = S
			default:
				cart = nil
			}
			if cart != nil {
				cart.id = id
				id++
				carts = append(carts, cart)
			}
		}
	}

	return track, carts
}
