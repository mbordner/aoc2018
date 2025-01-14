package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common/file"
	"math"
	"strings"
)

type Polymer struct {
	data []byte
}

func (p *Polymer) String() string {
	return string(p.data)
}

func (p *Polymer) Len() int {
	return len(p.data)
}

func (p *Polymer) Types() []byte {
	typesMap := make(map[byte]bool)
	for i := 0; i < p.Len(); i++ {
		t := p.data[i]
		if t > 'Z' {
			t -= byte(32)
		}
		typesMap[t] = true
	}
	types := make([]byte, 0, len(typesMap))
	for t := range typesMap {
		types = append(types, t)
	}
	return types
}

func (p *Polymer) Remove(t byte) *Polymer {
	T := t
	if t >= 'A' && t <= 'Z' {
		t += byte(32)
	} else {
		T -= byte(32)
	}
	nd := make([]byte, 0, len(p.data))
	for i := 0; i < len(p.data); i++ {
		if p.data[i] != T && p.data[i] != t {
			nd = append(nd, p.data[i])
		}
	}
	return NewPolymer(string(nd))
}

func (p *Polymer) React() int {
	removed := 0
	change := true
	for change {
		change = false
		if len(p.data) > 1 {
			for i := 1; i < len(p.data); i++ {
				if p.reacts(p.data[i], p.data[i-1]) {
					change = true
					if i == 1 {
						p.data = p.data[i+1:]
					} else {
						p.data = append(p.data[0:i-1], p.data[i+1:]...)
					}
					i--
					removed += 2
				}
			}
		}
	}

	return removed
}

func (p *Polymer) reacts(A, a byte) bool {
	if A > a {
		A, a = a, A
	}
	if A+byte(32) == a {
		return true
	}
	return false
}

func NewPolymer(data string) *Polymer {
	return &Polymer{data: []byte(data)}
}

func main() {
	/*
		p := NewPolymer(`aA`)
		fmt.Println(p.React(), p.Len(), p.String())
		p = NewPolymer(`abBA`)
		fmt.Println(p.React(), p.Len(), p.String())
		p = NewPolymer(`abAB`)
		fmt.Println(p.React(), p.Len(), p.String())
		p = NewPolymer(`aabAAB`)
		fmt.Println(p.React(), p.Len(), p.String())
		p = NewPolymer(`dabAcCaCBAcCcaDA`)
		fmt.Println(p.React(), p.Len(), p.String())
	*/

	p := getPolymer("../data.txt")

	var bestP *Polymer
	var bestT byte
	minLength := math.MaxUint32

	typesP := p.Types()
	for _, t := range typesP {
		op := p.Remove(t)
		op.React()
		orl := op.Len()
		if orl < minLength {
			minLength = orl
			bestP = op
			bestT = t
		}
	}

	fmt.Println(string(bestT), minLength, bestP)
}

func getPolymer(filename string) *Polymer {
	content, _ := file.GetContent(filename)
	return NewPolymer(strings.TrimSpace(string(content)))
}
