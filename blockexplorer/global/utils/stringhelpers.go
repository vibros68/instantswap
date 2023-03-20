package utils

import (
	"strconv"
	"strings"
)

// GetStringBefore returns string before x string value
func GetStringBefore(value string, a string) string {
	// Get substring before a string.
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}
	return value[0:pos]
}

func StrToFloat(str string) float64 {
	f, _ := strconv.ParseFloat(str, 64)
	return f
}

func StrToInt(str string) int {
	num, _ := strconv.Atoi(str)
	return num
}
