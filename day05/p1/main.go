package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common/file"
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
	fmt.Println(p.React(), p.Len(), p.String())
}

func getPolymer(filename string) *Polymer {
	content, _ := file.GetContent(filename)
	return NewPolymer(strings.TrimSpace(string(content)))
}
