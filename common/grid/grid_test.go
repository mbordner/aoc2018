package grid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CloneCondenseExpand(t *testing.T) {
	g := ExpandGrid([]byte{'0', '1', '2', '/', '3', '4', '5', '/', '6', '7', '8'}, '/')
	gc := g.Condense('/')
	o := g.Clone()
	oc := o.Condense('/')
	assert.Equal(t, gc, oc)
}

func Test_NewCondense(t *testing.T) {
	g := NewGrid(2, 2, 'x')
	gc := g.Condense('/')
	assert.Equal(t, "xx/xx", string(gc))
}

func Test_Rotate(t *testing.T) {
	g := ExpandGrid([]byte(`012/345/678`), '/')
	o := g.RotateRight()
	oc := o.Condense('/')
	assert.Equal(t, "630/741/852", string(oc))

	g = ExpandGrid([]byte(`01/23`), '/')
	o = g.RotateRight()
	oc = o.Condense('/')
	assert.Equal(t, "20/31", string(oc))
}

func Test_FlipHorizontal(t *testing.T) {
	g := ExpandGrid([]byte(`012/345/678`), '/')
	o := g.FlipHorizontal()
	oc := o.Condense('/')
	assert.Equal(t, "678/345/012", string(oc))

	g = ExpandGrid([]byte(`01/23`), '/')
	o = g.FlipHorizontal()
	oc = o.Condense('/')
	assert.Equal(t, "23/01", string(oc))
}

func Test_FlipVertical(t *testing.T) {
	g := ExpandGrid([]byte(`012/345/678`), '/')
	o := g.FlipVertical()
	oc := o.Condense('/')
	assert.Equal(t, "210/543/876", string(oc))

	g = ExpandGrid([]byte(`01/23`), '/')
	o = g.FlipVertical()
	oc = o.Condense('/')
	assert.Equal(t, "10/32", string(oc))
}

func Test_Read(t *testing.T) {
	g := ExpandGrid([]byte(`012/345/678`), '/')
	o := g.Read(0, 0, 1, 1)
	oc := o.Condense('/')
	assert.Equal(t, "01/34", string(oc))

	o = g.Read(0, 0, 2, 2)
	oc = o.Condense('/')
	p := g.Clone()
	pc := p.Condense('/')
	assert.Equal(t, oc, pc)
}

func Test_Write(t *testing.T) {
	g := ExpandGrid([]byte(`012/345/678`), '/')
	o := ExpandGrid([]byte(`ab/cd`), '/')
	g.Write(1, 1, o)
	gc := g.Condense('/')
	assert.Equal(t, "012/3ab/6cd", string(gc))
}
