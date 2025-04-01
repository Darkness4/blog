// Package math are utilities for math.
package math

func MaxI(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func MinI(a int, b int) int {
	if b > a {
		return a
	}
	return b
}
