// Package mat provides simple generic math functions
package mat

import "golang.org/x/exp/constraints"

// Min returns the minimum value of the given two values
func Min[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

// Max returns the maximum value of the given two values
func Max[T constraints.Ordered](x, y T) T {
	if x > y {
		return x
	}
	return y
}
