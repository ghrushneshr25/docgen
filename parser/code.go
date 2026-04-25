package parser

import (
	"os"
	"regexp"
	"strings"
)

func ExtractCode(file string) string {
	data, _ := os.ReadFile(file)
	re := regexp.MustCompile(`(?s)/\*.*?\*/`)
	out := string(re.ReplaceAll(data, []byte("")))
	out = strings.TrimSpace(trimLeadingMetadataComments(out))
	return FormatGoSource(out)
}

// trimLeadingMetadataComments removes file-header lines like "// @problem: ..." that
// are used for docgen but should not appear in the rendered solution block.
func trimLeadingMetadataComments(s string) string {
	lines := strings.Split(s, "\n")
	i := 0
	for i < len(lines) {
		t := strings.TrimSpace(lines[i])
		if t == "" {
			i++
			continue
		}
		if isMetadataCommentLine(t) {
			i++
			continue
		}
		break
	}
	return strings.Join(lines[i:], "\n")
}

func isMetadataCommentLine(trimmed string) bool {
	if !strings.HasPrefix(trimmed, "//") {
		return false
	}
	rest := strings.TrimSpace(trimmed[2:])
	return strings.HasPrefix(rest, "@")
}
