package cmath

import (
	"math"
	"sort"
)

type Number interface {
	int | int32 | int64 | float32 | float64
}

type IntNumber interface {
	int | int32 | int64 | uint64
}

func Factorial[V IntNumber](v V) V {
	if v == V(1) {
		return v
	}
	return v * Factorial[V](v-V(1))
}

var (
	MaxInt   = int(^uint(0) >> 1)
	MaxInt64 = int64(^uint64(0) >> 1)
)

func Factors[V IntNumber](v V) []V {
	var factors []V

	// Check from 1 to the square root of n
	for i := V(1); i <= V(math.Sqrt(float64(v))); i++ {
		if v%i == 0 {
			factors = append(factors, i)

			// If the divisor is not the square root, add the corresponding divisor
			if v/i != i {
				factors = append(factors, v/i)
			}
		}
	}

	sort.Slice(factors, func(i, j int) bool {
		return factors[i] < factors[j]
	})

	return factors
}

func IsPrime[V IntNumber](v V) bool {
	if v <= 1 {
		return false
	} else if v == 2 {
		return true
	} else if v%2 == 0 {
		return false
	}
	sqrt := V(math.Sqrt(float64(v)))
	for i := V(3); i <= sqrt; i += 2 {
		if v%i == 0 {
			return false
		}
	}
	return true
}

func Sum[V Number](vs []V) V {
	sum := V(0)
	for _, v := range vs {
		sum += v
	}
	return sum
}

// Sums returns all positive pairs from 0...v that add up to v sorted
func Sums[V IntNumber](v V) []V {
	vs := make([]V, 0, v*2)
	for i, j := V(0), v; i <= j; i, j = i+1, j-1 {
		vs = append(vs, []V{i, j}...)
	}
	return vs
}

func Product[V Number](vs []V) V {
	product := vs[0]
	for _, v := range vs[1:] {
		product *= v
	}
	return product
}

func PrimeFactors[V IntNumber](v V) []V {
	var factors []V

	// Check for divisibility by 2
	for v%2 == 0 {
		factors = append(factors, 2)
		v /= 2
	}

	// Check for divisibility by odd numbers
	for i := V(3); i*i <= v; i += 2 {
		for v%i == 0 {
			factors = append(factors, i)
			v /= i
		}
	}

	// If n is a prime number greater than 2
	if v > 2 {
		factors = append(factors, v)
	}

	return factors
}

/*
*

combinations formula:
C(n,k)  =  n! /   k! (n-k)!
n choose 5 (order doesn't matter)

Permutations:
P(n,k) = n! / (n-k)!

P(n,k) from n card deck, how many k card decks can be made.

k! - how many different ways to arrange k cards
*/
