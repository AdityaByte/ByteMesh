package utils

import "strings"

func CheckEmptyField(field string) bool {
	field = strings.TrimSpace(field)
	if field == "" {
		return true
	}
	return false
}