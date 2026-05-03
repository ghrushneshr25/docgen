package renderer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"docgen/parser"
	"docgen/utils"
)

// RenderReadme writes readme.md at repo root: per category, two subsections when both
// exist — **Concepts** table first, then **Problems** table. Each list is ordered by
// @index (then title-prefix fallback) via parser.OrderedSources.
func RenderReadme(codeDir, categoryOrderFile, outPath string) error {
	entries, err := os.ReadDir(codeDir)
	if err != nil {
		return err
	}
	entries = utils.OrderedDirEntries(entries, categoryOrderFile)

	var b strings.Builder
	b.WriteString("# 📚 DSA Index\n\n")

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		catName := e.Name()
		pattern := filepath.Join(codeDir, catName, "*.go")
		files, err := filepath.Glob(pattern)
		if err != nil {
			return err
		}
		var conceptPaths, problemPaths []string
		for _, f := range files {
			if utils.IsTestFile(f) {
				continue
			}
			meta := parser.ParseMetadata(f)
			if parser.DocType(meta) == "concept" {
				conceptPaths = append(conceptPaths, f)
			} else {
				problemPaths = append(problemPaths, f)
			}
		}
		if len(conceptPaths) == 0 && len(problemPaths) == 0 {
			continue
		}

		concepts := parser.OrderedSources(conceptPaths)
		problems := parser.OrderedSources(problemPaths)

		header := utils.FormatTitle(catName)
		b.WriteString("## 📂 ")
		b.WriteString(header)
		b.WriteString("\n\n")

		if len(concepts) > 0 {
			b.WriteString("### Concepts\n\n")
			writeReadmeTableHeader(&b, "Concept")
			writeReadmeTableRows(&b, codeDir, catName, concepts)
			b.WriteByte('\n')
		}
		if len(problems) > 0 {
			b.WriteString("### Problems\n\n")
			writeReadmeTableHeader(&b, "Problem")
			writeReadmeTableRows(&b, codeDir, catName, problems)
			b.WriteByte('\n')
		}
		b.WriteByte('\n')
	}

	return os.WriteFile(outPath, []byte(b.String()), 0o644)
}

func writeReadmeTableHeader(b *strings.Builder, firstCol string) {
	b.WriteString("| ")
	b.WriteString(firstCol)
	b.WriteString(" | Difficulty | Tags | Description | Code | Tests |\n")
	b.WriteString("|--------|------------|------|-------------|------|-------|\n")
}

func writeReadmeTableRows(b *strings.Builder, codeDir, catName string, sources []parser.OrderedSource) {
	for _, src := range sources {
		file := src.Path
		meta := src.Meta
		stem := strings.TrimSuffix(filepath.Base(file), ".go")
		display := meta["problem"]
		if display == "" {
			display = meta["title"]
		}
		if display == "" {
			display = stem
		}
		difficulty := meta["difficulty"]
		if difficulty == "" {
			difficulty = "—"
		}
		tags := meta["tags"]
		if tags == "" {
			tags = "—"
		}
		desc := parser.ExtractDescriptionSection(file)
		if desc == "" {
			desc = "—"
		}

		display = sanitizeTableCell(display)
		difficulty = sanitizeTableCell(difficulty)
		tags = sanitizeTableCell(tags)
		desc = sanitizeTableCell(desc)

		relCode := filepath.ToSlash(filepath.Join(catName, filepath.Base(file)))
		codeCell := fmt.Sprintf("[Link](./%s)", relCode)

		testPath := filepath.Join(codeDir, catName, "test", stem+"_test.go")
		testCell := "—"
		if _, err := os.Stat(testPath); err == nil {
			relTest := filepath.ToSlash(filepath.Join(catName, "test", stem+"_test.go"))
			testCell = fmt.Sprintf("[Link](./%s)", relTest)
		}

		b.WriteString("| ")
		b.WriteString(display)
		b.WriteString(" | ")
		b.WriteString(difficulty)
		b.WriteString(" | ")
		b.WriteString(tags)
		b.WriteString(" | ")
		b.WriteString(desc)
		b.WriteString(" | ")
		b.WriteString(codeCell)
		b.WriteString(" | ")
		b.WriteString(testCell)
		b.WriteString(" |\n")
	}
}

func sanitizeTableCell(s string) string {
	s = strings.ReplaceAll(s, "|", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.TrimSpace(s)
}
