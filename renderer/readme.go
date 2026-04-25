package renderer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"docgen/parser"
	"docgen/utils"
)

// RenderReadme writes readme.md at repo root in the same shape as generate_readme.sh
// (per-category tables: Problem, Difficulty, Tags, Description, Code, Tests).
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
		var list []string
		for _, f := range files {
			if utils.IsTestFile(f) {
				continue
			}
			list = append(list, f)
		}
		if len(list) == 0 {
			continue
		}
		utils.SortSourcePathsByBase(list)

		header := utils.FormatTitle(catName)
		b.WriteString("## 📂 ")
		b.WriteString(header)
		b.WriteString("\n\n")
		b.WriteString("| Problem | Difficulty | Tags | Description | Code | Tests |\n")
		b.WriteString("|--------|------------|------|-------------|------|-------|\n")

		for _, file := range list {
			meta := parser.ParseMetadata(file)
			stem := strings.TrimSuffix(filepath.Base(file), ".go")
			problemName := meta["problem"]
			if problemName == "" {
				problemName = meta["title"]
			}
			if problemName == "" {
				problemName = stem
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

			problemName = sanitizeTableCell(problemName)
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
			b.WriteString(problemName)
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
		b.WriteByte('\n')
	}

	return os.WriteFile(outPath, []byte(b.String()), 0o644)
}

func sanitizeTableCell(s string) string {
	s = strings.ReplaceAll(s, "|", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.TrimSpace(s)
}
