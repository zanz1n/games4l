package telemetry

import (
	"unicode"

	"github.com/go-playground/validator/v10"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	normalizer = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	validate   = validator.New()
)

func eliminateDuplicates(arr []TelemetryUnit) []TelemetryUnit {
	cache := make(map[string]struct{})

	newArr := []TelemetryUnit{}

	var (
		ok bool
		v  TelemetryUnit
	)
	for _, v = range arr {
		if _, ok = cache[v.ID.Hex()]; !ok {
			newArr = append(newArr, v)
			cache[v.ID.Hex()] = struct{}{}
		}
	}

	return newArr
}
