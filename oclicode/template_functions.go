package oclicode

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bcdxn/opencli/oclispec"
)

var breakRE = regexp.MustCompile(`(?:[^\S\r\n]|[^0-9A-Za-z])+`)

func funcmap() map[string]any {
	return map[string]any{
		"PascalCase":   pascalCase,
		"CamelCase":    camelCase,
		"EscapeString": escapeString,
		"Inc":          increment,
		"ToString":     toString,
		"Dec":          decrement,
		"AddCmdDepth":  addCmdDepth,
	}
}

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

func decrement(n int) int {
	return n - 1
}

func addCmdDepth(n int, cmd oclispec.Command) int {
	cmdSegments := regexp.MustCompile(`[^\S\r\n]`).Split(cmd.Name, -1)
	// we subtract 1 because we don't want to include the binary in the command length
	return len(cmdSegments) + n - 1
}
