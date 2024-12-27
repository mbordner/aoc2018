package datastructure

import (
	"cmp"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type SomeData []int

func TestNewAnyHeap(t *testing.T) {

	h := NewAnyHeap[SomeData](func(a, b SomeData) int {
		return cmp.Compare(a[1], b[1])
	})

	assert.NotNil(t, h)

	h.Unshift(SomeData{1, 2})

	assert.Equal(t, `[1 2]`, fmt.Sprintf(`%v`, h.Peek()))

	h.Unshift(SomeData{1, 3})

	assert.Equal(t, `[1 2]`, fmt.Sprintf(`%v`, h.Peek()))

	assert.Equal(t, 2, h.Len())

	h.Unshift(SomeData{0, 4})

	assert.Equal(t, `[1 2]`, fmt.Sprintf(`%v`, h.Peek()))

	assert.Equal(t, 3, h.Len())

	h.Unshift(SomeData{0, 0})

	assert.Equal(t, `[0 0]`, fmt.Sprintf(`%v`, h.Peek()))

	assert.Equal(t, 4, h.Len())

	sd := h.Shift()

	assert.Equal(t, 3, h.Len())

	assert.Equal(t, `[0 0]`, fmt.Sprintf(`%v`, sd))

	assert.Equal(t, `[1 2]`, fmt.Sprintf(`%v`, h.Peek()))

	assert.Equal(t, `[1 3]`, fmt.Sprintf(`%v`, h.Get(1)))
	assert.Equal(t, `[0 4]`, fmt.Sprintf(`%v`, h.Get(2)))
}
