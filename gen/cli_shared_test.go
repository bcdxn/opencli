package gen

import "testing"

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"find-by-status", "FindByStatus"},
		{"photo_urls", "PhotoUrls"},
		{"petstore", "Petstore"},
		{"a-b-c-d", "ABCD"},
		{"AlreadyPascal", "AlreadyPascal"},
		{"", ""},
	}
	for _, tt := range tests {
		got := toPascalCase(tt.in)
		if got != tt.want {
			t.Errorf("toPascalCase(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"path-to-req-body", "pathToReqBody"},
		{"first_name", "firstName"},
		{"firstName", "firstName"},
		{"SingleWord", "singleWord"},
		{"", ""},
	}
	for _, tt := range tests {
		got := toCamelCase(tt.in)
		if got != tt.want {
			t.Errorf("toCamelCase(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestToGoPackageName(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"find-by-status", "findbystatus"},
		{"UploadImage", "uploadimage"},
		{"hello_world!", "helloworld"},
	}
	for _, tt := range tests {
		got := toGoPackageName(tt.in)
		if got != tt.want {
			t.Errorf("toGoPackageName(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestToGoType(t *testing.T) {
	tests := []struct {
		typ      string
		variadic bool
		want     string
	}{
		{"string", false, "string"},
		{"integer", false, "int64"},
		{"boolean", false, "bool"},
		{"number", false, "float64"},
		{"unknown", false, "string"},
		{"string", true, "[]string"},
		{"integer", true, "[]int64"},
	}
	for _, tt := range tests {
		got := toGoType(tt.typ, tt.variadic)
		if got != tt.want {
			t.Errorf("toGoType(%q, %v) = %q, want %q", tt.typ, tt.variadic, got, tt.want)
		}
	}
}

func TestToTSType(t *testing.T) {
	tests := []struct {
		typ      string
		variadic bool
		want     string
	}{
		{"string", false, "string"},
		{"integer", false, "number"},
		{"boolean", false, "boolean"},
		{"number", false, "number"},
		{"unknown", false, "string"},
		{"string", true, "string[]"},
		{"integer", true, "number[]"},
	}
	for _, tt := range tests {
		got := toTSType(tt.typ, tt.variadic)
		if got != tt.want {
			t.Errorf("toTSType(%q, %v) = %q, want %q", tt.typ, tt.variadic, got, tt.want)
		}
	}
}

func TestBuildMethodName(t *testing.T) {
	got := buildMethodName([]string{"pet", "add"})
	if got != "PetAdd" {
		t.Errorf("buildMethodName = %q, want %q", got, "PetAdd")
	}
}

func TestSplitAliases(t *testing.T) {
	shorthand, extra := splitAliases([]string{"v", "verbose", "V"})
	if shorthand != "v" {
		t.Errorf("shorthand = %q, want %q", shorthand, "v")
	}
	if len(extra) != 2 || extra[0] != "verbose" || extra[1] != "V" {
		t.Errorf("extra = %v, want [verbose V]", extra)
	}

	// No single-char alias
	shorthand, extra = splitAliases([]string{"long1", "long2"})
	if shorthand != "" || len(extra) != 2 {
		t.Errorf("no shorthand case failed: shorthand=%q, extra=%v", shorthand, extra)
	}
}
