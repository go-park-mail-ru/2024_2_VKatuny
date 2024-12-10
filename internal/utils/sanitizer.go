package utils

import (
	"github.com/microcosm-cc/bluemonday"
)

func SanitizeString(input string) string {
	p := bluemonday.UGCPolicy()
	return p.Sanitize(input)
}
