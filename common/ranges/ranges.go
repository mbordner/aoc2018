package ranges

import (
	"github.com/pkg/errors"
	"sort"
)

type Number interface {
	int | int32 | float32 | int64 | float64
}

type side int

const (
	leftSide side = iota
	rightSide
)

type el[T Number] struct {
	val T
	s   side
}

type elements[T Number] []el[T]

func (e elements[T]) Len() int {
	return len(e)
}

func (e elements[T]) Less(a, b int) bool {
	if e[a].val < e[b].val {
		return true
	}
	if e[a].val == e[b].val {
		if e[a].s == leftSide && e[b].s == rightSide {
			return true
		}
	}
	return false
}

func (e elements[T]) Swap(a, b int) {
	e[a], e[b] = e[b], e[a]
}

type Collection[T Number] struct {
	values []T
}

func (c *Collection[T]) ValuePairs() []T {
	return c.values
}

func (c *Collection[T]) Len() T {
	var l T
	for i := 0; i < len(c.values); i += 2 {
		l += c.values[i+1] - c.values[i] + 1
	}
	return l
}

func (c *Collection[T]) Add(l, r T) ([]T, error) {

	if r < l {
		return nil, errors.New("right value cannot be less than left value")
	}

	els := make(elements[T], 0, len(c.values)+2)

	if len(c.values) > 0 {
		for i, v := range c.values {
			s := leftSide
			if i%2 == 1 {
				s = rightSide
			}
			els = append(els, el[T]{val: v, s: s})
		}
	}

	els = append(els, el[T]{val: l, s: leftSide})
	els = append(els, el[T]{val: r, s: rightSide})

	sort.Sort(els)

	values := make([]T, 0, len(c.values)+2)

	leftValStack := make([]T, 0, 10)

	for i := 0; i < len(els); i++ {
		if els[i].s == leftSide {
			leftValStack = append(leftValStack, els[i].val)
		} else if els[i].s == rightSide {
			leftVal := leftValStack[len(leftValStack)-1]
			leftValStack = leftValStack[0 : len(leftValStack)-1]
			if len(leftValStack) == 0 {
				rightVal := els[i].val
				if len(values) > 0 {
					if leftVal == values[len(values)-1] {
						values[len(values)-1] = rightVal
						continue
					}
				}
				values = append(values, leftVal)
				values = append(values, rightVal)
			}
		}

	}

	c.values = values

	return values, nil
}
