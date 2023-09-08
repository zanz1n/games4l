package utils

import "unsafe"

func B2S(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func S2B(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
