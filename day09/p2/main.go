package main

import (
	"fmt"
)

type IntNumber interface {
	int | int32 | int64
}

type Node[T IntNumber] struct {
	val  T
	prev *Node[T]
	next *Node[T]
}

func abs[T IntNumber](x T) T {
	if x >= 0 {
		return x
	}
	return -x
}

func (n *Node[T]) Get(away T) *Node[T] {
	cur := n
	for i := T(0); i < abs(away); i++ {
		if away > 0 {
			cur = cur.next
		} else {
			cur = cur.prev
		}
	}
	return cur
}

type Circle[T IntNumber] struct {
	ptr *Node[T]
}

func NewCircle[T IntNumber]() *Circle[T] {
	c := new(Circle[T])
	c.ptr = &Node[T]{val: T(0)}
	c.ptr.prev = c.ptr
	c.ptr.next = c.ptr
	return c
}

type Values[T IntNumber] []T

func (vs Values[T]) Sum() T {
	var sum T
	for _, v := range vs {
		sum += v
	}
	return sum
}

func (c *Circle[T]) Add(val T) Values[T] {
	var nPtr *Node[T]
	var values Values[T]
	if val%23 == 0 {
		nPtr = c.ptr.Get(-6)

		rmVal := nPtr.prev

		nPtr.prev = rmVal.prev
		rmVal.prev.next = nPtr

		values = Values[T]{val, rmVal.val}
	} else {
		nPtr = c.ptr.Get(1)
		n := &Node[T]{val: val}

		n.next = nPtr.next
		n.next.prev = n

		nPtr.next = n
		n.prev = nPtr

		nPtr = n
	}
	c.ptr = nPtr
	return values
}

func main() {
	players := make([]int64, 423)
	ptr := 0

	circle := NewCircle[int64]()

	n := int64(1)
	for n <= int64(7194400) {
		vals := circle.Add(n)
		players[ptr] += vals.Sum()
		ptr++
		if ptr == len(players) {
			ptr = 0
		}
		n++
	}

	maxScore := int64(0)
	for _, v := range players {
		if v > maxScore {
			maxScore = v
		}
	}

	fmt.Println(maxScore)
}
