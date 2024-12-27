package datastructure

import (
	"container/heap"
	"github.com/pkg/errors"
)

type AnyHeap[T any] struct {
	data  []T
	cmp   func(a, b T) int
	index int
}

func (h *AnyHeap[T]) Less(i, j int) bool { return h.cmp(h.data[i], h.data[j]) < 0 }
func (h *AnyHeap[T]) Swap(i, j int)      { h.data[i], h.data[j] = h.data[j], h.data[i] }

func (h *AnyHeap[T]) Push(x any) {
	h.data = append(h.data, x.(T))
}

func (h *AnyHeap[T]) Pop() any {
	l := len(h.data)
	v := h.data[l-1]
	h.data = h.data[0 : l-1]
	return v
}

// do not use the above methods

func (h *AnyHeap[T]) Len() int { return len(h.data) }

func (h *AnyHeap[T]) Shift() T {
	return heap.Pop(h).(T)
}

func (h *AnyHeap[T]) Unshift(v T) {
	heap.Push(h, v)
}

func (h *AnyHeap[T]) Get(i int) T {
	return h.data[i]
}

func (h *AnyHeap[T]) Peek() T {
	return h.Get(0)
}

func (h *AnyHeap[T]) RewindNext() {
	h.index = 0
}

func (h *AnyHeap[T]) RewindTo(t T) {
	for i := 0; i < len(h.data); i++ {
		if h.cmp(h.data[i], t) == 0 {
			h.index = i
			break
		}
	}
}

func (h *AnyHeap[T]) HasNext() bool {
	if h.index < h.Len() {
		return true
	}
	return false
}

func (h *AnyHeap[T]) Next() (T, error) {
	if h.HasNext() {
		v := h.Get(h.index)
		h.index++
		return v, nil
	}
	var t T
	return t, errors.New("exhausted next")
}

func (h *AnyHeap[T]) PeekNext() (T, error) {
	if h.HasNext() {
		v := h.Get(h.index)
		return v, nil
	}
	var t T
	return t, errors.New("exhausted next")
}

func (h *AnyHeap[T]) PeekAfterNext() (T, error) {
	if h.index+1 < h.Len() {
		v := h.Get(h.index + 1)
		return v, nil
	}
	var t T
	return t, errors.New("exhausted next")
}

// NewAnyHeap cmp function return 0 if a == b, < 0 if a < b, > 0 if a > b
func NewAnyHeap[T any](cmp func(a, b T) int) *AnyHeap[T] {
	h := new(AnyHeap[T])
	h.cmp = cmp
	h.data = make([]T, 0)
	heap.Init(h)
	return h
}
