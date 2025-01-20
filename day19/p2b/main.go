package main

import "fmt"

func main() {

	num := 10551288

	numFactorsSum := 0

	r4 := 1
	r1 := 1

	/*
		[1 2 3 4 6 8 11 12 17 22 24 33 34 44 51 66 68 88 102 132 136 187 204 264 374 408 561 748 1122 1496 2244 2351 4488 4702 7053 9404 14106 18808 25861 28212 39967 51722 56424 77583 79934 103444 119901 155166 159868 206888 239802 310332 319736 439637 479604 620664 879274 959208 1318911 1758548 2637822 3517096 5275644 10551288]
		30481920
		[2 2 2 3 11 17 2351]
		10551288
		false
	*/

	lastR1 := num

	for r4 <= lastR1 { // if (r4 <= r1), it will be the sum of prime factors + 1 (sum includes 1)
		r1 = lastR1
		for r1 > r4 {
			if r4*r1 == num {
				numFactorsSum += r4 + r1
				lastR1 = r1
			}
			r1--
		}
		r4++
	}

	fmt.Println(numFactorsSum)

}
