// Package ptr are utilities for pointer.
package ptr

func Ref[T any](v T) *T {
	return &v
}
