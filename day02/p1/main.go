package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common/file"
	"strings"
)

func main() {
	ss := getData("../data.txt")

	var c2, c3 int

	for _, s := range ss {
		var has2, has3 bool
		cs := counts(s)
		for _, c := range cs {
			if c == 2 {
				has2 = true
			} else if c == 3 {
				has3 = true
			}
		}
		if has2 {
			c2++
		}
		if has3 {
			c3++
		}
	}

	fmt.Println(c2 * c3)
}

func counts(s string) map[byte]int {
	cs := make(map[byte]int)

	for _, b := range []byte(s) {
		if c, e := cs[b]; e {
			cs[b] = c + 1
		} else {
			cs[b] = 1
		}
	}

	return cs
}

func getData(filename string) []string {
	lines, _ := file.GetLines(filename)
	ss := make([]string, len(lines))
	for i, line := range lines {
		tokens := strings.Fields(line)
		ss[i] = tokens[0]
	}
	return ss
}
