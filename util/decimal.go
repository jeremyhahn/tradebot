package util

import "strings"

func ParseDecimalPlace(value string) int32 {
	pieces := strings.Split(value, ".")
	if len(pieces) < 2 {
		return 0
	}
	idx := strings.Index(pieces[1], "1")
	if idx > -1 {
		return int32(idx + 1)
	}
	return 0
}
