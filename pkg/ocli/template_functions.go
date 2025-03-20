package ocli

import (
	"regexp"
	"strings"
)

var breakRE = regexp.MustCompile(`(?:[^\S\r\n]|[^0-9A-Za-z])+`)

func pascalCase(s string) string {
	segments := breakRE.Split(s, -1)
	pascalCasedSegments := []string{}
	for i := range segments {
		if len(segments[i]) > 0 {
			pascalCasedSegments = append(pascalCasedSegments, strings.ToUpper(string(segments[i][0]))+strings.ToLower(segments[i][1:]))
		}
	}

	return strings.Join(pascalCasedSegments, "")
}
