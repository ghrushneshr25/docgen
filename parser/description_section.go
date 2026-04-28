package parser

import (
	"os"
	"strings"
)

// ExtractDescriptionSection returns plain text from the first block comment between
// @section: Description and the next @section: line (same behavior as generate_readme.sh).
func ExtractDescriptionSection(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}
	content := string(data)
	start := strings.Index(content, "/*")
	end := strings.LastIndex(content, "*/")
	if start == -1 || end == -1 || end <= start {
		return ""
	}
	block := content[start+2 : end]
	lines := strings.Split(block, "\n")
	inDesc := false
	var parts []string
	for _, raw := range lines {
		lineTrim := strings.TrimSpace(raw)
		if strings.HasPrefix(lineTrim, "@section:") {
			name := strings.TrimSpace(strings.TrimPrefix(lineTrim, "@section:"))
			if strings.EqualFold(name, "Description") {
				inDesc = true
				continue
			}
			if inDesc {
				break
			}
			continue
		}
		if !inDesc {
			continue
		}
		t := strings.TrimSpace(raw)
		t = strings.TrimPrefix(t, "*")
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		parts = append(parts, t)
	}
	out := strings.Join(parts, " ")
	out = strings.Join(strings.Fields(out), " ")
	return out
}
