package ocli

import (
	"fmt"
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

func camelCase(s string) string {
	segments := breakRE.Split(s, -1)
	pascalCasedSegments := []string{strings.ToLower(segments[0])}
	for i := 1; i < len(segments); i++ {
		if len(segments[i]) > 0 {
			pascalCasedSegments = append(pascalCasedSegments, strings.ToUpper(string(segments[i][0]))+strings.ToLower(segments[i][1:]))
		}
	}

	return strings.Join(pascalCasedSegments, "")
}

var newLineRE = regexp.MustCompile(`\r?\n`)
var doubleQuoteRE = regexp.MustCompile(`"`)

func escapeString(s string) string {
	escapedStr := doubleQuoteRE.ReplaceAllString(s, `\"`)
	return newLineRE.ReplaceAllString(escapedStr, `\n`)
}

func increment(n int) int {
	return n + 1
}

func toString(v any) string {
	if v == nil {
		return ""
	}

	return fmt.Sprintf("%v", v)
}
