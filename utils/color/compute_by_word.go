// Package color are utilities for color.
package color

import (
	"crypto/sha256"
	"fmt"
	"math"
)

// ComputeByWord returns a CSS string color based on the word.
//
// It return --label-r, --label-g, --label-b, --label-h, --label-s, --label-l based on the hash of the word.
func ComputeByWord(word string) string {
	r, g, b := RGBFromWord(word)
	h, s, l := HSLFromRGB(r, g, b)
	return fmt.Sprintf(
		"--label-r: %d; --label-g: %d; --label-b: %d; --label-h: %d; --label-s: %d; --label-l: %d;",
		r,
		g,
		b,
		h,
		s,
		l,
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

// HSLFromRGB returns the HSL values based on the RGB values.
func HSLFromRGB(r, g, b int) (hue, saturation, lightness int) {
	// Convert RGB to normalized values (0-1).
	rn := float64(r) / 255.0
	gn := float64(g) / 255.0
	bn := float64(b) / 255.0

	// Find the maximum and minimum values to compute lightness.
	M := math.Max(math.Max(rn, gn), bn)
	m := math.Min(math.Min(rn, gn), bn)
	l := (M + m) / 2

	// Initialize variables for hue and saturation.
	var h, s float64

	// Check for saturation.
	if M == m {
		// Achromatic case: no hue.
		h = 0
		s = 0
	} else {
		// Chromatic case.
		d := M - m
		if l > 0.5 {
			s = d / (2 - M - m)
		} else {
			s = d / (M + m)
		}

		// Compute hue based on the channel with the maximum value.
		switch M {
		case rn:
			h = (gn - bn) / d
			if gn < bn {
				h += 6
			}
		case gn:
			h = ((bn - rn) / d) + 2
		case bn:
			h = ((rn - gn) / d) + 4
		}

		h *= 60 // Convert to degrees.
	}

	// Ensure the values are within the valid range.
	hue = int(math.Round(h))
	if hue < 0 {
		hue += 360
	}
	saturation = int(math.Round(s * 100))
	lightness = int(math.Round(l * 100))

	return hue, saturation, lightness
}
