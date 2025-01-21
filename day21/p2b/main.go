package main

import "fmt"

func main() {

	r3s := make(map[int]bool)
	r4s := make(map[int]bool)
	lastR3 := 0

	r3 := 2176960
	r4 := 65536

	for {
		r3 += r4 & 255

		r3 &= 16777215
		r3 *= 65899
		r3 &= 16777215

		if 256 > r4 {
			if _, e := r3s[r3]; !e {
				r3s[r3] = true
				lastR3 = r3
			}
			r4 = r3 | 65536
			if _, e := r4s[r4]; !e {
				r4s[r4] = true
			} else {
				break
			}
			r3 = 2176960
		} else {
			r4 /= 256
		}
	}

	fmt.Println(lastR3)

}
