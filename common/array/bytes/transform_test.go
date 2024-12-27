package bytes_test

import (
	"aoc2021/common/array/bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRotate(t *testing.T) {

	a1 := [][]byte{
		{'1', '2', '3'},
		{'4', '5', '6'},
		{'7', '8', '9'},
	}

	e1 := [][]byte{
		{'7', '4', '1'},
		{'8', '5', '2'},
		{'9', '6', '3'},
	}
	a1key := fmt.Sprintf("%v", a1)
	e1key := fmt.Sprintf("%v", e1)

	r1 := bytes.Rotate(a1)
	r1key := fmt.Sprintf("%v", r1)
	assert.Equal(t, e1key, r1key)

	r1 = bytes.Rotate(r1)
	r1 = bytes.Rotate(r1)
	r1 = bytes.Rotate(r1)

	r1key = fmt.Sprintf("%v", r1)
	assert.Equal(t, a1key, r1key)

	a2 := [][]byte{
		{'1', '2', '3', '4'},
		{'5', '6', '7', '8'},
		{'9', 'A', 'B', 'C'},
		{'D', 'E', 'F', 'G'},
	}

	e2 := [][]byte{
		{'D', '9', '5', '1'},
		{'E', 'A', '6', '2'},
		{'F', 'B', '7', '3'},
		{'G', 'C', '8', '4'},
	}
	e2key := fmt.Sprintf("%v", e2)

	r2 := bytes.Rotate(a2)
	r2key := fmt.Sprintf("%v", r2)
	assert.Equal(t, e2key, r2key)
}

func TestFlip(t *testing.T) {
	a1 := [][]byte{
		{'1', '2', '3'},
		{'4', '5', '6'},
		{'7', '8', '9'},
	}

	r1 := bytes.Flip(bytes.Horizontal, a1)

	assert.EqualValues(t, [][]byte{
		{'3', '2', '1'},
		{'6', '5', '4'},
		{'9', '8', '7'},
	}, r1)

	r2 := bytes.Flip(bytes.Vertical, a1)

	assert.EqualValues(t, [][]byte{
		{'7', '8', '9'},
		{'4', '5', '6'},
		{'1', '2', '3'},
	}, r2)
}
