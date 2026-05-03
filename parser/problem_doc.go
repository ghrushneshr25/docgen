package parser

import (
	"os"
	"regexp"
	"strings"
)

// docSegment is one @section or @subsection block in the leading /* */ doc.
type docSegment struct {
	Name  string
	Lines []string
}

var reApproachSection = regexp.MustCompile(`(?i)^approach\s+\d+`)

// isApproachSection is true for titles like "Approach 1 - Brute Force".
func isApproachSection(name string) bool {
	return reApproachSection.MatchString(strings.TrimSpace(name))
}

func isSummarySection(name string) bool {
	n := strings.ToLower(strings.TrimSpace(name))
	return n == "summary" || n == "comparing approaches"
}

// splitDocBlockIntoSegments parses the leading comment block into ordered segments.
func splitDocBlockIntoSegments(file string) ([]docSegment, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	content := string(data)
	start := strings.Index(content, "/*")
	end := strings.LastIndex(content, "*/")
	if start == -1 || end == -1 {
		return nil, nil
	}
	block := content[start+2 : end]
	rawLines := strings.Split(block, "\n")

	var segs []docSegment
	var buf []string
	sectionName := ""

	flush := func() {
		if sectionName == "" {
			buf = buf[:0]
			return
		}
		if len(buf) == 0 {
			segs = append(segs, docSegment{Name: sectionName, Lines: nil})
			buf = buf[:0]
			return
		}
		lines := append([]string(nil), buf...)
		segs = append(segs, docSegment{Name: sectionName, Lines: lines})
		buf = buf[:0]
	}

	for _, rawLine := range rawLines {
		lineTrim := strings.TrimSpace(rawLine)
		if strings.HasPrefix(lineTrim, "@section:") {
			flush()
			sectionName = strings.TrimSpace(strings.TrimPrefix(lineTrim, "@section:"))
			continue
		}
		if strings.HasPrefix(lineTrim, "@subsection:") {
			flush()
			sectionName = strings.TrimSpace(strings.TrimPrefix(lineTrim, "@subsection:"))
			continue
		}
		buf = append(buf, rawLine)
	}
	flush()
	return segs, nil
}

// segmentBodyMarkdown renders one segment's body (no ## heading).
func segmentBodyMarkdown(sectionName string, lines []string) string {
	if isAlgorithmSection(sectionName) {
		trimmed := trimBlankRuns(lines)
		trimmed = dedentLines(trimmed)
		var b strings.Builder
		b.WriteString("```text\n")
		for _, ln := range trimmed {
			b.WriteString(ln)
			b.WriteString("\n")
		}
		b.WriteString("```\n\n")
		return b.String()
	}
	var b strings.Builder
	for _, ln := range lines {
		t := strings.TrimSpace(ln)
		if t == "" {
			b.WriteString("\n")
			continue
		}
		b.WriteString(t)
		b.WriteString("\n")
	}
	b.WriteString("\n")
	return b.String()
}

// segmentToMarkdown renders "## Name" plus body (same heading level as generate-binary-strings.mdx).
func segmentToMarkdown(seg docSegment) string {
	var b strings.Builder
	b.WriteString("\n## ")
	b.WriteString(seg.Name)
	b.WriteString("\n\n")
	b.WriteString(segmentBodyMarkdown(seg.Name, seg.Lines))
	return b.String()
}

// ProblemApproachDoc is one approach group for multi-approach problem pages (prose only;
// Go snippets are attached in main via ExtractOperationCodeOrdered).
type ProblemApproachDoc struct {
	Markdown string // "## Approach …" + "## Algorithm" + "## Notes" + …
}

// ParseProblemMultiApproach builds preamble, per-approach markdown, and optional summary
// when the doc uses one or more @section: Approach N … blocks with matching // @operation
// functions in file order.
func ParseProblemMultiApproach(file string) (preamble string, approaches []ProblemApproachDoc, summary string, ok bool) {
	segs, err := splitDocBlockIntoSegments(file)
	if err != nil || len(segs) == 0 {
		return "", nil, "", false
	}

	var preambleSegs []docSegment
	var groups [][]docSegment
	var summarySegs []docSegment

	var curGroup []docSegment
	inApproaches := false

	for _, s := range segs {
		if isSummarySection(s.Name) {
			if len(curGroup) > 0 {
				groups = append(groups, curGroup)
				curGroup = nil
			}
			summarySegs = append(summarySegs, s)
			inApproaches = false
			continue
		}
		if isApproachSection(s.Name) {
			if len(curGroup) > 0 {
				groups = append(groups, curGroup)
			}
			curGroup = []docSegment{s}
			inApproaches = true
			continue
		}
		if inApproaches {
			curGroup = append(curGroup, s)
			continue
		}
		preambleSegs = append(preambleSegs, s)
	}
	if len(curGroup) > 0 {
		groups = append(groups, curGroup)
	}

	// Multi-approach layout only when there are 2+ approaches (matches kth-node style).
	// A single @section: Approach 1 … with one @operation stays on the legacy template
	// (Description → Algorithm → Notes → one Solution), like check-if-array-is-sorted.
	if len(groups) < 2 {
		return "", nil, "", false
	}

	n, err := OperationDirectiveCount(file)
	if err != nil || n != len(groups) {
		return "", nil, "", false
	}

	var pb strings.Builder
	for _, s := range preambleSegs {
		pb.WriteString(segmentToMarkdown(s))
	}
	preamble = strings.TrimSpace(pb.String())

	for _, g := range groups {
		var ab strings.Builder
		for _, s := range g {
			ab.WriteString(segmentToMarkdown(s))
		}
		approaches = append(approaches, ProblemApproachDoc{
			Markdown: strings.TrimSpace(ab.String()),
		})
	}

	if len(summarySegs) > 0 {
		var sb strings.Builder
		for _, s := range summarySegs {
			sb.WriteString(segmentToMarkdown(s))
		}
		summary = strings.TrimSpace(sb.String())
	}

	return preamble, approaches, summary, true
}
