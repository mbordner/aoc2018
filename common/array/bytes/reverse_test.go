package bytes_test

import (
	"aoc2021/common/array/bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReverse(t *testing.T) {
	a := []byte{'1', '2', '3'}
	b := bytes.Reverse(a)
	assert.Equal(t, a[2], b[0])
	assert.Equal(t, a[1], b[1])
	assert.Equal(t, a[0], b[2])

	a = []byte{'1', '2', '3', '4'}
	b = bytes.Reverse(a)
	assert.Equal(t, a[3], b[0])
	assert.Equal(t, a[2], b[1])
	assert.Equal(t, a[1], b[2])
	assert.Equal(t, a[0], b[3])
}
