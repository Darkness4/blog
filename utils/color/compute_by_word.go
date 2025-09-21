// Package color are utilities for color.
package color

import (
	"crypto/sha256"
	"fmt"
	"math"
)

// ComputeByWord returns a CSS string color based on the word.
//
// It return --label-color based on the hash of the word.
func ComputeByWord(word string) string {
	r, g, b := RGBFromWord(word)
	return fmt.Sprintf(
		"--label-color: rgb(%d, %d, %d);", r, g, b,
	)
}

// RGBFromWord returns the RGB values based on the word.
func RGBFromWord(word string) (r, g, b int) {
	// Calculate SHA256 hash of the word.
	hash := sha256.Sum256([]byte(word))

	// Use first 4 bytes of the hash to generate RGB values.
	r = int(hash[0])
	g = int(hash[1])
	b = int(hash[2])

	// Scale RGB values to the range of 0-255.
	r = int(math.Round(float64(r) / 255.0 * 128))
	g = int(math.Round(float64(g) / 255.0 * 128))
	b = int(math.Round(float64(b) / 255.0 * 128))

	// Return the RGB values.
	return r, g, b
}
