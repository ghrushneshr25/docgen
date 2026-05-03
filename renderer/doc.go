package renderer

import "docgen/parser"

// OperationSubsection is one @subsection under Operations: prose and Go for that block.
type OperationSubsection struct {
	Title string
	Prose string
	Code  string
}

// ProblemApproach is one approach block for multi-approach problem docs (Algorithm + Notes + Go).
type ProblemApproach struct {
	Markdown string // includes ## Approach N …, ## Algorithm, ## Notes, …
	Code     string // gofmt'd snippet for this approach
}

type Doc struct {
	Title    string
	// SidebarPosition is 1-based order in the category sidebar (matches sort by
	// leading digits in Title, then title text).
	SidebarPosition int
	Type            string
	Sections string // problems: full ParseSections markdown (single-solution layout)
	Code     string // problems: full implementation; concepts: unused (use OperationsCode)
	// When non-empty, problem.tpl renders per-approach Algorithm / Notes / Solution (Go) like generate-binary-strings.mdx × N.
	ProblemPreamble   string
	ProblemApproaches []ProblemApproach
	ProblemSummary    string
	Subtests []parser.SubtestInfo // each t.Run name, @desc, and code
	HasTests bool                 // true if any subtest should be shown
	Meta     map[string]string
	Structs  []parser.StructInfo

	// concepts only (fixed render order: Description → Structure → Operations)
	ConceptDescription    string
	ConceptStructureIntro string // @section: Structure before first @subsection
	ConceptOperationsMD   string // prose under @section: Operations
	OperationsCode        string // gofmt’d functions only (joined) when no OperationSubsections
	// when non-empty, interleave ### title, prose, and code (from banner sections in the .go file)
	OperationSubsections []OperationSubsection

	Output string
}
