package utils

import "unicode"

func FirstUpper(str string) string {
	if len(str) <= 0 {
		return str
	}
	rns := []rune(str)
	rns[0] = unicode.ToUpper(rns[0])
	return string(rns)
}

func SliceContains[Type comparable](s []Type, item Type) bool {
	for _, ai := range s {
		if ai == item {
			return true
		}
	}

	return false
}
