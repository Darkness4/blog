package math

import (
	"fmt"
	"strings"
)

func FormatNumber(num float64) string {
	var formatted string
	switch {
	case num >= 1e9:
		formatted = fmt.Sprintf("%.1fb", num/1e9)
	case num >= 1e6:
		formatted = fmt.Sprintf("%.1fm", num/1e6)
	case num >= 1e3:
		formatted = fmt.Sprintf("%.1fk", num/1e3)
	default:
		formatted = fmt.Sprintf("%.1f", num)
	}

	return strings.TrimSuffix(formatted, ".0")
}
