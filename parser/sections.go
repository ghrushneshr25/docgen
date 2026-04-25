package parser

import (
	"os"
	"strings"
)

func ParseSections(file string) string {
	data, _ := os.ReadFile(file)
	content := string(data)

	start := strings.Index(content, "/*")
	end := strings.LastIndex(content, "*/")

	if start == -1 || end == -1 {
		return ""
	}

	block := content[start+2 : end]
	rawLines := strings.Split(block, "\n")

	var out strings.Builder
	var buf []string
	sectionName := ""

	flushBuf := func() {
		if len(buf) == 0 {
			return
		}
		if isAlgorithmSection(sectionName) {
			trimmed := trimBlankRuns(buf)
			trimmed = dedentLines(trimmed)
			out.WriteString("```text\n")
			for _, ln := range trimmed {
				out.WriteString(ln)
				out.WriteString("\n")
			}
			out.WriteString("```\n\n")
		} else {
			for _, ln := range buf {
				t := strings.TrimSpace(ln)
				if t == "" {
					out.WriteString("\n")
					continue
				}
				out.WriteString(t)
				out.WriteString("\n")
			}
			out.WriteString("\n")
		}
		buf = buf[:0]
	}

	for _, rawLine := range rawLines {
		lineTrim := strings.TrimSpace(rawLine)
		if strings.HasPrefix(lineTrim, "@section:") {
			flushBuf()
			sectionName = strings.TrimSpace(strings.TrimPrefix(lineTrim, "@section:"))
			out.WriteString("\n## ")
			out.WriteString(sectionName)
			out.WriteString("\n\n")
			continue
		}
		if strings.HasPrefix(lineTrim, "@subsection:") {
			flushBuf()
			sectionName = strings.TrimSpace(strings.TrimPrefix(lineTrim, "@subsection:"))
			out.WriteString("\n### ")
			out.WriteString(sectionName)
			out.WriteString("\n\n")
			continue
		}
		buf = append(buf, rawLine)
	}
	flushBuf()

	return out.String()
}

func isAlgorithmSection(name string) bool {
	return strings.EqualFold(strings.TrimSpace(name), "algorithm")
}

func trimBlankRuns(lines []string) []string {
	i, j := 0, len(lines)-1
	for i <= j && strings.TrimSpace(lines[i]) == "" {
		i++
	}
	for j >= i && strings.TrimSpace(lines[j]) == "" {
		j--
	}
	if i > j {
		return nil
	}
	return lines[i : j+1]
}

func dedentLines(lines []string) []string {
	min := -1
	for _, ln := range lines {
		if strings.TrimSpace(ln) == "" {
			continue
		}
		n := leadingIndentLen(ln)
		if min < 0 || n < min {
			min = n
		}
	}
	if min <= 0 {
		out := make([]string, len(lines))
		for i, ln := range lines {
			out[i] = strings.TrimRight(ln, " \t\r")
		}
		return out
	}
	out := make([]string, len(lines))
	for i, ln := range lines {
		if strings.TrimSpace(ln) == "" {
			out[i] = ""
			continue
		}
		if len(ln) >= min {
			out[i] = strings.TrimRight(ln[min:], " \t\r")
		} else {
			out[i] = strings.TrimRight(ln, " \t\r")
		}
	}
	return out
}

func leadingIndentLen(s string) int {
	n := 0
	for _, r := range s {
		if r == ' ' || r == '\t' {
			n++
			continue
		}
		break
	}
	return n
}
