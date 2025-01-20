package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common/cmath"
)

func main() {
	num := 10551288
	fmt.Println(cmath.Factors(num))
	fmt.Println(cmath.Sum(cmath.Factors(num)))
	fmt.Println(cmath.PrimeFactors(num))
	fmt.Println(cmath.Product(cmath.PrimeFactors(num)))
	fmt.Println(cmath.Sum(cmath.PrimeFactors(num)))
	fmt.Println(cmath.IsPrime(num))
}
