package bytes

type Pos struct {
	R int
	C int
}

func FindMasked(bitmap [][]byte, mask [][]byte) []Pos {
	positions := make([]Pos, 0, 10)

	height := len(mask)
	width := len(mask[0])

	for j := 0; j <= len(bitmap)-height; j++ {
		for i := 0; i <= len(bitmap[0])-width; i++ {
			pos := Pos{R: j, C: i}
			if CheckMask(bitmap, mask, pos) {
				positions = append(positions, pos)
			}
		}
	}

	return positions
}

func CheckMask(bitmap [][]byte, mask [][]byte, pos Pos) bool {
	height := len(mask)
	width := len(mask[0])
	// j,i position within bitmaps
	// n,m position within mask
	for j, n := pos.R, 0; j < pos.R+height; j, n = j+1, n+1 {
		for i, m := pos.C, 0; i < pos.C+width; i, m = i+1, m+1 {
			if mask[n][m] != 0 {
				if bitmap[j][i] != mask[n][m] {
					return false
				}
			}
		}
	}
	return true
}

func ApplyMask(bitmap [][]byte, mask [][]byte, pos Pos) {
	height := len(mask)
	width := len(mask[0])
	// j,i position within bitmaps
	// n,m position within mask
	for j, n := pos.R, 0; j < pos.R+height; j, n = j+1, n+1 {
		for i, m := pos.C, 0; i < pos.C+width; i, m = i+1, m+1 {
			if mask[n][m] != 0 {
				bitmap[j][i] ^= mask[n][m]
			}
		}
	}
}
