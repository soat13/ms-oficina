package string_helper

import (
	"regexp"
	"strconv"
)

func OnlyNumbers(s string) string {
	re := regexp.MustCompile(`[^0-9]`)
	return re.ReplaceAllString(s, "")
}

func StringToIntOrDefault(s string, def int) int {
	if s == "" {
		return def
	}
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return def
}
