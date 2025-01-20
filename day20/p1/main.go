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

func (rm RoomMap) resolve(cur *Room, pattern string) *Room {
	if len(pattern) == 0 || cur == nil {
		return cur
	}
	if pattern[0] == ParenStart {
		closeIndex := rm.searchCloseParen(0, pattern)
		subPattern := pattern[1:closeIndex]
		continuePattern := pattern[closeIndex+1:]
		options := rm.splitOptions(subPattern)
		for _, option := range options {
			last := rm.resolve(cur, option)
			rm.resolve(last, continuePattern)
		}
		return nil
	} else {
		var np common.Pos
		var dir common.Pos
		switch pattern[0] {
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
		nr = rm.Has(np)
		if nr == nil {
			nr = &Room{id: np}
			rm[nr.id] = nr
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
			cur.e = cur
		}
		return rm.resolve(nr, pattern[1:])
	}
	panic("unreachable")
}

func (rm RoomMap) splitOptions(pattern string) []string {
	var options []string
	optionStart := 0
	for i := 0; i < len(pattern); i++ {
		if pattern[i] == ParenStart {
			i = rm.searchCloseParen(i, pattern)
		} else if pattern[i] == Pipe {
			option := pattern[optionStart:i]
			options = append(options, option)
			optionStart = i + 1 // it's possible that there could be || in the pattern, but it doesn't exist in data..
		}
	}
	options = append(options, pattern[optionStart:])
	return options
}

type cpMem struct {
	i int
	p string
}

var (
	mem = make(map[cpMem]int)
)

func (rm RoomMap) searchCloseParen(start int, pattern string) int {
	sm := cpMem{i: start, p: pattern}
	if i, e := mem[sm]; e {
		return i
	}
	ps := make(common.Stack[int], 0, 10)
	if pattern[start] == ParenStart {
		ps.Push(start)
		for i := start + 1; i < len(pattern); i++ {
			if pattern[i] == ParenStart {
				ps.Push(i)
			} else if pattern[i] == ParenEnd {
				if ps.Len() == 1 {
					mem[sm] = i
					return i
				} else {
					ps.Pop()
				}
			}
		}
	}
	return start
}

func (rm RoomMap) Resolve(pattern string) {
	cur := &Room{id: common.Pos{}}
	start := strings.Index(pattern, "^")
	end := strings.Index(pattern, "$")
	pattern = pattern[start+1 : end]
	rm[cur.id] = cur
	rm.resolve(cur, pattern)
}

func main() {
	content, _ := file.GetContent("../data.txt")
	rm := make(RoomMap)
	rm.Resolve(strings.TrimSpace(string(content)))
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

	fmt.Println(maxLen)
	fmt.Println(maxPath)

}
