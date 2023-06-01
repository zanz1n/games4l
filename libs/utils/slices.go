package utils

func SliceContains[Type comparable](s []Type, item Type) bool {
	for _, ai := range s {
		if ai == item {
			return true
		}
	}

	return false
}
