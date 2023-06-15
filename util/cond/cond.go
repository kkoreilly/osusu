// Package cond provides functions for executing inline conditional statements
package cond

// If returns ifTrue if the given condition is true; otherwise it returns the zero value of T
func If[T any](cond bool, ifTrue T) T {
	if cond {
		return ifTrue
	}
	var res T
	return res
}

// IfElse returns ifTrue if the given condition is true; otherwise, it returns ifFalse
func IfElse[T any](cond bool, ifTrue T, ifFalse T) T {
	if cond {
		return ifTrue
	}
	return ifFalse
}
