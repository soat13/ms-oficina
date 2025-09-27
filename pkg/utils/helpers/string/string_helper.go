package string_helper

import "regexp"

func OnlyNumbers(s string) string {
	re := regexp.MustCompile(`[^0-9]`)
	return re.ReplaceAllString(s, "")
}
