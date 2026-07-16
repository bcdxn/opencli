package gen

import "strings"

// specArgEntry holds the minimal spec data for an argument (used in getSpec*Cmd()).
type specArgEntry struct {
	Name    string
	Summary string
}

// specFlagEntry holds the minimal spec data for a flag (used in getSpec*Cmd()).
type specFlagEntry struct {
	Name    string
	Summary string
	Aliases []string
}

// commandFileCoreTmplData is the shared template data used across framework command files.
// Framework-specific command data structs should embed this and only add unique fields.
type commandFileCoreTmplData struct {
	ModuleVersion            string
	Binary                   string
	BinaryPascal             string
	FuncName                 string
	SpecFuncName             string
	OutPath                  string
	Segment                  string
	IsRoot                   bool
	IsGroup                  bool
	IsHidden                 bool
	MethodName               string
	ArgsTypeName             string
	FlagsTypeName            string
	Summary                  string
	Description              string
	Aliases                  []string
	CommandLine              string
	VisibleChildren          bool
	VisibleArgs              bool
	VisibleFlags             bool
	SpecArgs                 []specArgEntry
	SpecFlags                []specFlagEntry
	CommandModifiers         []string
	ArgsModifiers            []string
	FlagsModifiers           []string
	PassthroughArgsModifiers []string
}

func appendSegment(parentSegments []string, segment string) []string {
	segments := make([]string, len(parentSegments)+1)
	copy(segments, parentSegments)
	segments[len(parentSegments)] = segment
	return segments
}

func buildCommandFileCore(
	cmdSegment string,
	summary string,
	description string,
	aliases []string,
	hidden bool,
	visibleChildren bool,
	visibleArgs bool,
	visibleFlags bool,
	commandModifiers []string,
	argsModifiers []string,
	flagsModifiers []string,
	passthroughModifiers []string,
	segments []string,
	binary string,
	binaryPascal string,
	moduleVersion string,
	isRoot bool,
	isGroup bool,
) commandFileCoreTmplData {
	methodName := buildMethodName(segments)

	return commandFileCoreTmplData{
		ModuleVersion:            moduleVersion,
		Binary:                   binary,
		BinaryPascal:             binaryPascal,
		Segment:                  cmdSegment,
		IsRoot:                   isRoot,
		IsGroup:                  isGroup,
		IsHidden:                 hidden,
		MethodName:               methodName,
		ArgsTypeName:             methodName + "Args",
		FlagsTypeName:            methodName + "Flags",
		Summary:                  summary,
		Description:              strings.TrimRight(description, "\n"),
		Aliases:                  aliases,
		VisibleChildren:          visibleChildren,
		VisibleArgs:              visibleArgs,
		VisibleFlags:             visibleFlags,
		CommandModifiers:         commandModifiers,
		ArgsModifiers:            argsModifiers,
		FlagsModifiers:           flagsModifiers,
		PassthroughArgsModifiers: passthroughModifiers,
	}
}
