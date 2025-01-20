package main

import "fmt"

func main() {

	num := 10551288

	numFactorsSum := 0

	r4 := 1
	r1 := 1

	/*
		factors of num:
		[1 2 3 4 6 8 11 12 17 22 24 33 34 44 51 66 68 88 102 132 136 187 204 264 374 408 561 748 1122 1496 2244 2351 4488 4702 7053 9404 14106 18808 25861 28212 39967 51722 56424 77583 79934 103444 119901 155166 159868 206888 239802 310332 319736 439637 479604 620664 879274 959208 1318911 1758548 2637822 3517096 5275644 10551288]
		sum of the factors: 30481920

		prime factors of num:
		[2 2 2 3 11 17 2351]
		sum of prime factors: 2388
		product of multiplying them all back together: 10551288
	*/

	lastR1 := num

	for r4 < lastR1 {
		r1 = lastR1
		for r1 > r4 {
			if r4*r1 == num {
				fmt.Println(r4, r1)
				numFactorsSum += r4 + r1
				lastR1 = r1
			}
			r1--
		}
		r4++
	}

	fmt.Println(numFactorsSum)

}
