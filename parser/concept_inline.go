package parser

import (
	"os"
	"regexp"
	"strings"
)

// ConceptInlineBody is prose and operation titles from // @structure / // @operation lines
// in the .go file (in-file docs). Complements the leading /* */ block when @subsections
// are omitted and prose lives next to the code.
type ConceptInlineBody struct {
	// StructDoc maps Go type name (e.g. ListNode) to prose from the preceding
	// // @description line. Keys come from the `type` declaration, not the // @structure label
	// (e.g. "List Node" vs type ListNode).
	StructDoc map[string]string
	// Operations in source order; Title from // @operation <Title>
	Operations []OperationProse
}

var reInlineDescription = regexp.MustCompile(`^\s*//\s*@description\s*(.*)$`)
var reInlineTime = regexp.MustCompile(`^\s*//\s*@time:\s*(.*)$`)
var reInlineSpace = regexp.MustCompile(`^\s*//\s*@space:\s*(.*)$`)
var reTypeStruct = regexp.MustCompile(`^\s*type\s+(\w+)\s+struct\s*`)

// ParseConceptInlineFromFile reads // @structure, // @description, // @operation, // @time, and // @space
// in line comments. Structure prose is taken from // @description; operation prose is those lines
// up to the following func declaration, rendered as description + Time + Space in markdown.
func ParseConceptInlineFromFile(path string) (ConceptInlineBody, error) {
	out := ConceptInlineBody{StructDoc: make(map[string]string)}
	data, err := os.ReadFile(path)
	if err != nil {
		return out, err
	}
	lines := strings.Split(string(data), "\n")
	for i := 0; i < len(lines); {
		line := lines[i]
		if opTitle, ok := titleFromOperationTagLine(line); ok {
			prose, next := parseOpMetadataLines(lines, i+1)
			out.Operations = append(out.Operations, OperationProse{Title: opTitle, Prose: strings.TrimSpace(prose)})
			i = next
			continue
		}
		if _, ok := titleFromStructureTagLine(line); ok {
			typName, st := parseStructureAfterTag(lines, i+1)
			if typName != "" && st.ok {
				out.StructDoc[typName] = st.firstDesc
			}
			// st.j = index after the `type` line. If we did not advance, skip the @structure line.
			if st.j > i+1 {
				i = st.j
			} else {
				i++
			}
			continue
		}
		i++
	}
	return out, nil
}

// isGoFuncLine reports whether the line is the start of a func declaration in Go.
func isGoFuncLine(s string) bool {
	s = strings.TrimSpace(s)
	if len(s) < 4 || s[:4] != "func" {
		return false
	}
	if len(s) == 4 {
		return true
	}
	r := s[4]
	return r == ' ' || r == '\t' || r == '('
}

// parseOpMetadataLines returns markdown prose and the line index of the func keyword.
func parseOpMetadataLines(lines []string, start int) (prose string, funcLine int) {
	var b strings.Builder
	for j := start; j < len(lines); j++ {
		ln := lines[j]
		if isGoFuncLine(ln) {
			return b.String(), j
		}
		if m := reInlineDescription.FindStringSubmatch(ln); len(m) == 2 {
			if t := strings.TrimSpace(m[1]); t != "" {
				if b.Len() > 0 {
					// Tight: hard line break (not a blank paragraph) between parts
					b.WriteString("  \n")
				}
				b.WriteString(t)
			}
		} else if m := reInlineTime.FindStringSubmatch(ln); len(m) == 2 {
			if b.Len() > 0 {
				b.WriteString("  \n")
			}
			b.WriteString("Time: " + strings.TrimSpace(m[1]))
		} else if m := reInlineSpace.FindStringSubmatch(ln); len(m) == 2 {
			if b.Len() > 0 {
				b.WriteString("  \n")
			}
			b.WriteString("Space: " + strings.TrimSpace(m[1]))
		}
	}
	return b.String(), len(lines)
}

type nextStructParse struct {
	j         int
	firstDesc string
	ok        bool
}

// parseStructureAfterTag reads optional // @description, skips blanks, then the next `type Name struct` line.
func parseStructureAfterTag(lines []string, start int) (typName string, next nextStructParse) {
	j := start
	for j < len(lines) && strings.TrimSpace(lines[j]) == "" {
		j++
	}
	if j < len(lines) {
		if m := reInlineDescription.FindStringSubmatch(lines[j]); m != nil {
			if t := strings.TrimSpace(m[1]); t != "" {
				next.firstDesc = t
				next.ok = true
			}
			j++
		}
	}
	for j < len(lines) && strings.TrimSpace(lines[j]) == "" {
		j++
	}
	if j < len(lines) {
		if m := reTypeStruct.FindStringSubmatch(lines[j]); len(m) == 2 {
			typName = m[1]
			j++
		}
	}
	next.j = j
	return typName, next
}
