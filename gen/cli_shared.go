package gen

import (
	"strings"
	"unicode"
)

// toPascalCase converts a kebab-case, snake_case, or camelCase string to PascalCase.
// "find-by-status" -> "FindByStatus", "photoUrls" -> "PhotoUrls", "petstore" -> "Petstore"
func toPascalCase(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '-' || r == '_'
	})
	var result strings.Builder
	for _, p := range parts {
		if len(p) == 0 {
			continue
		}
		runes := []rune(p)
		result.WriteRune(unicode.ToUpper(runes[0]))
		for _, r := range runes[1:] {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// toCamelCase converts kebab-case or snake_case to camelCase (yargs argv field convention).
// "path-to-req-body" -> "pathToReqBody"
func toCamelCase(s string) string {
	if s == "" {
		return s
	}

	// If the name is already a single token (e.g. "firstName"), preserve internal
	// casing and only lowercase the first rune.
	if !strings.ContainsAny(s, "-_") {
		r := []rune(s)
		r[0] = unicode.ToLower(r[0])
		return string(r)
	}

	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '-' || r == '_'
	})
	if len(parts) == 0 {
		return s
	}

	var result strings.Builder
	result.WriteString(strings.ToLower(parts[0]))
	for _, p := range parts[1:] {
		if len(p) == 0 {
			continue
		}
		runes := []rune(p)
		runes[0] = unicode.ToUpper(runes[0])
		result.WriteString(string(runes))
	}
	return result.String()
}

// toGoPackageName converts a command segment to a valid, lowercase Go package name.
// "find-by-status" -> "findbystatus", "upload-image" -> "uploadimage"
func toGoPackageName(s string) string {
	var result strings.Builder
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// toGoType maps an OpenCLI spec type string to a Go type.
func toGoType(t string, variadic bool) string {
	base := map[string]string{
		"string":  "string",
		"integer": "int64",
		"boolean": "bool",
		"number":  "float64",
	}[t]
	if base == "" {
		base = "string"
	}
	if variadic {
		return "[]" + base
	}
	return base
}

// toTSType maps an OpenCLI spec type to a TypeScript type string.
func toTSType(t string, variadic bool) string {
	base := map[string]string{
		"string":  "string",
		"integer": "number",
		"boolean": "boolean",
		"number":  "number",
	}[t]
	if base == "" {
		base = "string"
	}
	if variadic {
		return base + "[]"
	}
	return base
}

// buildMethodName constructs the ActionsInterface method name for a command.
// e.g. binaryPascal="Petstore", segments=["pet","add"] -> "PetstorePetAdd"
func buildMethodName(segments []string) string {
	parts := make([]string, 0, len(segments)+1)
	for _, seg := range segments {
		parts = append(parts, toPascalCase(seg))
	}
	return strings.Join(parts, "")
}

// splitAliases returns the first single-char alias as shorthand and all other aliases.
func splitAliases(aliases []string) (string, []string) {
	shorthand := ""
	extraAliases := make([]string, 0, len(aliases))
	for _, a := range aliases {
		if len(a) == 1 && shorthand == "" {
			shorthand = a
			continue
		}
		extraAliases = append(extraAliases, a)
	}
	return shorthand, extraAliases
}
