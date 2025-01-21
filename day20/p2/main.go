package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common"
	"github.com/mbordner/aoc2018/common/file"
	"strings"
)

const (
	Pipe       = '|'
	ParenStart = '('
	ParenEnd   = ')'
)

type RoomMap map[common.Pos]*Room

func (rm RoomMap) Has(p common.Pos) *Room {
	if r, e := rm[p]; e {
		return r
	}
	return nil
}

type Room struct {
	id common.Pos
	n  *Room
	e  *Room
	s  *Room
	w  *Room
}

type RoomsResolver struct {
	rooms   RoomMap
	mem     map[int]int
	pattern string
}

func (rr *RoomsResolver) resolve(cur *Room, s, e int) []*Room {
	if s == e {
		return []*Room{cur}
	}
	if rr.pattern[s] == ParenStart {
		closeIndex := rr.searchCloseParen(s, e)
		options := rr.splitOptions(s+1, closeIndex)
		lastsMap := make(RoomMap)
		for i := 0; i < len(options); i += 2 {
			for _, last := range rr.resolve(cur, options[i], options[i+1]) {
				if last != nil {
					lastsMap[last.id] = last
				}
			}
		}
		if closeIndex+1 < e {
			for _, last := range lastsMap {
				rr.resolve(last, closeIndex+1, e)
			}
			return nil
		}
		lasts := make([]*Room, 0, len(lastsMap))
		for _, last := range lastsMap {
			lasts = append(lasts, last)
		}
		return lasts
	} else {
		var np common.Pos
		var dir common.Pos
		switch rr.pattern[s] {
		case 'N':
			dir = common.DN
			np = cur.id.Add(dir)
		case 'E':
			dir = common.DE
			np = cur.id.Add(dir)
		case 'S':
			dir = common.DS
			np = cur.id.Add(dir)
		case 'W':
			dir = common.DW
			np = cur.id.Add(dir)
		}
		var nr *Room
		nr = rr.rooms.Has(np)
		if nr == nil {
			nr = &Room{id: np}
			rr.rooms[nr.id] = nr
		}
		switch dir {
		case common.DN:
			cur.n = nr
			nr.s = cur
		case common.DE:
			cur.e = nr
			nr.w = cur
		case common.DS:
			cur.s = nr
			nr.n = cur
		case common.DW:
			cur.w = nr
			nr.e = cur
		}
		return rr.resolve(nr, s+1, e)
	}
	panic("unreachable")
}

func (rr *RoomsResolver) splitOptions(s, e int) []int {
	var options []int
	optionStart := s
	for i := s; i < e; i++ {
		if rr.pattern[i] == ParenStart {
			i = rr.searchCloseParen(i, e)
		} else if rr.pattern[i] == Pipe {
			options = append(options, []int{optionStart, i}...)
			optionStart = i + 1 // it's possible that there could be || in the pattern, but it doesn't exist in data.
		}
	}
	options = append(options, []int{optionStart, e}...)
	return options
}

func (rr *RoomsResolver) searchCloseParen(s, e int) int {
	if i, exists := rr.mem[s]; exists {
		return i
	}
	ps := make(common.Stack[int], 0, 10)
	if rr.pattern[s] == ParenStart {
		ps.Push(s)
		for i := s + 1; i < e; i++ {
			if rr.pattern[i] == ParenStart {
				ps.Push(i)
			} else if rr.pattern[i] == ParenEnd {
				if ps.Len() == 1 {
					rr.mem[s] = i
					return i
				} else {
					ps.Pop()
				}
			}
		}
	}
	return s
}

func (rr *RoomsResolver) Resolve() RoomMap {
	cur := &Room{id: common.Pos{}}
	start := strings.Index(rr.pattern, "^")
	end := strings.Index(rr.pattern, "$")
	rr.rooms[cur.id] = cur
	rr.resolve(cur, start, end)
	return rr.rooms
}

func NewRoomsResolver(pattern string) *RoomsResolver {
	return &RoomsResolver{pattern: pattern, mem: make(map[int]int), rooms: make(RoomMap)}
}

// 702 too low
func main() {
	content, _ := file.GetContent("../data.txt")
	rr := NewRoomsResolver(strings.TrimSpace(string(content)))
	//rr = NewRoomsResolver(`^WNE$`)                                     // 3
	//rr = NewRoomsResolver(`^ENWWW(NEEE|SSE(EE|N))$`)                   // 10
	//rr = NewRoomsResolver(`^ENNWSWW(NEWS|)SSSEEN(WNSE|)EE(SWEN|)NNN$`) // 18
	rm := rr.Resolve()
	visited := make(RoomMap)
	prev := make(common.PosLinker)

	start := common.Pos{}
	queue := make(common.Queue[*Room], 0, len(rm))
	queue.Enqueue(rm[start])
	visited[start] = rm[start]

	for !queue.Empty() {
		cur := *(queue.Dequeue())
		nrs := make([]*Room, 0, 4)
		if cur.n != nil {
			nrs = append(nrs, cur.n)
		}
		if cur.e != nil {
			nrs = append(nrs, cur.e)
		}
		if cur.s != nil {
			nrs = append(nrs, cur.s)
		}
		if cur.w != nil {
			nrs = append(nrs, cur.w)
		}
		for _, nr := range nrs {
			if visited.Has(nr.id) == nil {
				visited[nr.id] = nr
				prev[nr.id] = cur.id
				queue.Enqueue(nr)
			}
		}
	}

	maxLen := 0
	maxPath := common.Positions{}
	paths := make(map[common.Pos]common.Positions)

	for _, r := range rm {
		if r.id != start {
			path := common.Positions{r.id}
			for p := prev[r.id]; p != start; p = prev[p] {
				path = append(common.Positions{p}, path...)
			}
			paths[r.id] = path
			if len(path) > maxLen {
				maxLen = len(path)
				maxPath = path
			}
		}
	}

	fmt.Println("number of rooms:", len(rm))
	fmt.Println("number of paths:", len(paths))
	fmt.Println("max path length:", maxLen)
	fmt.Println("max path:", maxPath)

	count := 0
	for _, path := range paths {
		if len(path) >= 1000 {
			count++
		}
	}

	fmt.Printf("%d rooms have a path that pass through at least 1000 doors.\n", count)

}
