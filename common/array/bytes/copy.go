package bytes

func Copy2D(dest [][]byte, src [][]byte, dr int, dc int, sr int, sc int, w int, h int) {
	for j, i := dr, sr; j < dr+h; j, i = j+1, i+1 {
		copy(dest[j][dc:dc+w], src[i][sc:sc+w])
	}
}
