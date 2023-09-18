package repository

import (
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var normalizer = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

func eliminateDuplicates(arr []mongoRegistryData) []mongoRegistryData {
	cache := make(map[string]struct{})

	newArr := []mongoRegistryData{}

	var ok bool
	for _, v := range arr {
		if _, ok = cache[v.ID.Hex()]; !ok {
			newArr = append(newArr, v)
			cache[v.ID.Hex()] = struct{}{}
		}
	}

	return newArr
}
