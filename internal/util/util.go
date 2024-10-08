package util

import (
	"strconv"
	"strings"
)

func FormatFloat(num float64) string {
	numStr := strconv.FormatFloat(num, 'f', -1, 64)
	if !strings.Contains(numStr, ".") {
		return numStr + ".0"
	}

	// Remove trailing zeros and decimal point if necessary
	parts := strings.Split(numStr, ".")
	parts[1] = strings.TrimRight(parts[1], "0")
	if len(parts[1]) == 0 {
		parts[1] = ""
	}
	return strings.Join(parts, ".")
}
