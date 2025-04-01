// Package blog are utilities for blog.
package blog

import (
	"fmt"
	"strings"
	"time"
)

func ExtractDate(filename string) (time.Time, error) {
	// Split the filename by "-"
	parts := strings.Split(filename, "-")

	// Check if there are enough parts in the filename
	if len(parts) < 3 {
		return time.Time{}, fmt.Errorf("invalid filename format: %s", filename)
	}

	// Extract the date part (YYYY-MM-DD)
	dateStr := strings.Join(parts[:3], "-")

	// Parse the date string into a time.Time object
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}
