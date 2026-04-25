package renderer

import "docgen/parser"

type Doc struct {
	Title    string
	Type     string
	Sections string // problems: full ParseSections markdown
	Code     string // problems: full implementation; concepts: unused (use OperationsCode)
	Subtests []parser.SubtestInfo // each t.Run name, @desc, and code
	HasTests bool                 // true if any subtest should be shown
	Meta     map[string]string
	Structs  []parser.StructInfo

	// concepts only (fixed render order: Description → Structure → Operations)
	ConceptDescription    string
	ConceptStructureIntro string // @section: Structure before first @subsection
	ConceptOperationsMD   string // prose under @section: Operations
	OperationsCode        string // gofmt’d functions only (joined)

	Output string
}
