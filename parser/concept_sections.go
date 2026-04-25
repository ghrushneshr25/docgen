package parser

import (
	"os"
	"strings"
)

// ConceptDoc is parsed prose from a concept file's leading /* */ block.
type ConceptDoc struct {
	Description          string
	StructureIntro       string
	StructureSubsections map[string]string // @subsection title -> prose (keys match type names)
	Operations           string
}

// ParseConceptDocBlock parses the doc block. For @section: Structure, text before the first
// @subsection becomes StructureIntro; each @subsection body is stored by title (no duplicate
// ### headings) so the template can render ### Name, prose, then struct code once.
func ParseConceptDocBlock(path string) (*ConceptDoc, error) {
	out := &ConceptDoc{StructureSubsections: make(map[string]string)}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := string(data)
	start := strings.Index(content, "/*")
	end := strings.LastIndex(content, "*/")
	if start == -1 || end == -1 || end <= start {
		return out, nil
	}
	block := content[start+2 : end]
	rawLines := strings.Split(block, "\n")

	var descB, opsB strings.Builder
	var structIntro strings.Builder
	current := ""
	var buf []string
	lastHeading := ""

	// inside Structure: which @subsection we're filling ("" = intro before first @subsection)
	structSubCurrent := ""

	flushTo := func(target *strings.Builder) {
		if target == nil || len(buf) == 0 {
			buf = buf[:0]
			return
		}
		writeFormattedBody(target, buf, lastHeading)
		buf = buf[:0]
	}

	targetFor := func() *strings.Builder {
		switch current {
		case "description":
			return &descB
		case "structure":
			return nil // structure uses structIntro / map, not structB
		case "operations":
			return &opsB
		default:
			return nil
		}
	}

	flushStructureBuf := func() {
		if len(buf) == 0 {
			return
		}
		text := formattedBodyString(buf, lastHeading)
		buf = buf[:0]
		if text == "" {
			return
		}
		if structSubCurrent == "" {
			if structIntro.Len() > 0 {
				structIntro.WriteString("\n\n")
			}
			structIntro.WriteString(text)
		} else {
			out.StructureSubsections[structSubCurrent] = text
		}
	}

	for _, rawLine := range rawLines {
		lineTrim := strings.TrimSpace(rawLine)
		if strings.HasPrefix(lineTrim, "@section:") {
			if current == "structure" {
				flushStructureBuf()
				structSubCurrent = ""
			} else {
				flushTo(targetFor())
			}

			name := strings.TrimSpace(strings.TrimPrefix(lineTrim, "@section:"))
			key := conceptSectionKey(name)
			if key == "" {
				opsB.WriteString("\n## ")
				opsB.WriteString(name)
				opsB.WriteString("\n\n")
				current = "operations"
			} else {
				current = key
			}
			lastHeading = name
			continue
		}
		if strings.HasPrefix(lineTrim, "@subsection:") {
			sub := strings.TrimSpace(strings.TrimPrefix(lineTrim, "@subsection:"))
			if current == "structure" {
				flushStructureBuf()
				structSubCurrent = sub
				lastHeading = sub
				continue
			}
			flushTo(targetFor())
			tb := targetFor()
			if tb != nil {
				tb.WriteString("\n### ")
				tb.WriteString(sub)
				tb.WriteString("\n\n")
			}
			lastHeading = sub
			continue
		}
		buf = append(buf, rawLine)
	}

	if current == "structure" {
		flushStructureBuf()
	} else {
		flushTo(targetFor())
	}

	out.Description = strings.TrimSpace(descB.String())
	out.StructureIntro = strings.TrimSpace(structIntro.String())
	out.Operations = strings.TrimSpace(opsB.String())
	return out, nil
}

func writeFormattedBody(target *strings.Builder, buf []string, lastHeading string) {
	s := formattedBodyString(buf, lastHeading)
	if s == "" {
		return
	}
	target.WriteString(s)
	target.WriteString("\n\n")
}

func formattedBodyString(buf []string, lastHeading string) string {
	if len(buf) == 0 {
		return ""
	}
	if isAlgorithmSection(lastHeading) {
		trimmed := trimBlankRuns(buf)
		trimmed = dedentLines(trimmed)
		var b strings.Builder
		b.WriteString("```text\n")
		for _, ln := range trimmed {
			b.WriteString(ln)
			b.WriteString("\n")
		}
		b.WriteString("```\n\n")
		return strings.TrimSpace(b.String())
	}
	var b strings.Builder
	for _, ln := range buf {
		t := strings.TrimSpace(ln)
		if t == "" {
			b.WriteString("\n")
			continue
		}
		b.WriteString(t)
		b.WriteString("\n")
	}
	return strings.TrimSpace(b.String())
}

func conceptSectionKey(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "description":
		return "description"
	case "structure":
		return "structure"
	case "operations":
		return "operations"
	default:
		return ""
	}
}
