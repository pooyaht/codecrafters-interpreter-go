package util

import (
	"strconv"
	"strings"
)

func FormatFloat(num float64) string {
	defaultStr := strconv.FormatFloat(num, 'f', -1, 64)

	var numStr string
	if strings.Contains(defaultStr, "e") || strings.Contains(defaultStr, "E") {
		numStr = strconv.FormatFloat(num, 'f', 1, 64)
	} else {
		numStr = defaultStr
	}

	if !strings.Contains(numStr, ".") {
		return numStr + ".0"
	}

	parts := strings.Split(numStr, ".")
	parts[1] = strings.TrimRight(parts[1], "0")
	if len(parts[1]) == 0 {
		parts[1] = "0"
	}
	return strings.Join(parts, ".")
}
