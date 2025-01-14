package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common/file"
	"strconv"
	"strings"
)

type Node struct {
	children []*Node
	metadata []int
}

func (n *Node) metaSum() int {
	metaSum := 0

	if len(n.children) > 0 {
		for _, m := range n.metadata {
			c := m - 1
			if c >= 0 && c < len(n.children) {
				metaSum += n.children[c].metaSum()
			}
		}
	} else {
		for _, m := range n.metadata {
			metaSum += m
		}
	}

	return metaSum
}

func NewNode(index int, data []string) (*Node, int) {
	// <header><child nodes><metadata>
	// header : <num: # child nodes><num: # of metadata>
	numChildren := atoi(data[index])
	numMetadata := atoi(data[index+1])

	n := &Node{children: make([]*Node, numChildren), metadata: make([]int, numMetadata)}

	ptr := index + 2
	for c := 0; c < numChildren; c++ {
		cn, read := NewNode(ptr, data)
		ptr += read
		n.children[c] = cn
	}

	for m := 0; m < numMetadata; m, ptr = m+1, ptr+1 {
		n.metadata[m] = atoi(data[ptr])
	}

	return n, ptr - index
}

func main() {
	root := getData("../data.txt")
	fmt.Println(root.metaSum())
}

func getData(filename string) *Node {
	lines, _ := file.GetLines(filename)
	data := strings.Split(strings.TrimSpace(string(lines[0])), " ")

	root, read := NewNode(0, data)
	if read != len(data) {
		panic("fail")
	}

	return root
}

func atoi(s string) int {
	val, _ := strconv.ParseInt(s, 10, 64)
	return int(val)
}
