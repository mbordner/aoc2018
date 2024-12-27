package grid

type Grid[T comparable] [][]T

func (g Grid[T]) Condense(delim T) []T {
	size := len(g)*len(g[0]) + len(g) - 1
	a := make([]T, 0, size)
	for y := 0; y < len(g); y++ {
		if y > 0 {
			a = append(a, delim)
		}
		for x := 0; x < len(g[y]); x++ {
			a = append(a, g[y][x])
		}
	}
	return a
}

func (g Grid[T]) RotateRight() Grid[T] {
	o := make(Grid[T], len(g[0]))
	for y := 0; y < len(g); y++ {
		o[y] = make([]T, len(g[y]))
		for x, j := 0, len(g)-1; x < len(g[y]); x, j = x+1, j-1 {
			o[y][x] = g[j][y]
		}
	}
	return o
}

func (g Grid[T]) Clone() Grid[T] {
	o := make(Grid[T], len(g))
	for y := range o {
		o[y] = make([]T, len(g[y]))
		for x := range o[y] {
			o[y][x] = g[y][x]
		}
	}
	return o
}

// FlipHorizontal flips on horizontal axis
func (g Grid[T]) FlipHorizontal() Grid[T] {
	o := g.Clone()
	for x := range g[0] {
		// reverse the array
		for i, j, h := 0, len(g)-1, len(g)/2; i < h; i, j = i+1, j-1 {
			o[i][x], o[j][x] = o[j][x], o[i][x]
		}
	}
	return o
}

// FlipVertical flips on vertical axis
func (g Grid[T]) FlipVertical() Grid[T] {
	o := g.Clone()
	for y := range o {
		// reverse the array
		for i, j, h := 0, len(o[y])-1, len(o[y])/2; i < h; i, j = i+1, j-1 {
			o[y][i], o[y][j] = o[y][j], o[y][i]
		}
	}
	return o
}

// Read reads to new grid
func (g Grid[T]) Read(startX, startY, endX, endY int) Grid[T] {
	w := endX - startX + 1
	h := endY - startY + 1
	o := make(Grid[T], h)
	for y := range o {
		o[y] = make([]T, w)
		for x := range o[y] {
			o[y][x] = g[startY+y][startX+x]
		}
	}
	return o
}

// Write writes in place
func (g Grid[T]) Write(startX, startY int, o Grid[T]) {
	for y := range o {
		for x := range o[y] {
			g[startY+y][startX+x] = o[y][x]
		}
	}
}

func NewGrid[T comparable](w, h int, fill T) Grid[T] {
	o := make(Grid[T], h)
	for y := range o {
		o[y] = make([]T, w)
		for x := range o[y] {
			o[y][x] = fill
		}
	}
	return o
}

func ExpandGrid[T comparable](condensed []T, delim T) Grid[T] {
	return SplitSlice[T](condensed, delim)
}

func SplitSlice[T comparable](slice []T, element T) [][]T {
	var result [][]T
	var current []T

	for _, v := range slice {
		if v == element {
			if len(current) > 0 {
				result = append(result, current)
				current = []T{}
			}
		} else {
			current = append(current, v)
		}
	}

	if len(current) > 0 {
		result = append(result, current)
	}

	return result
}
