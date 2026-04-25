package parser

import (
	"go/format"
	"strings"
)

// FormatGoSource applies standard gofmt to complete Go source (e.g. a full file after stripping docs).
func FormatGoSource(src string) string {
	src = strings.TrimSpace(src)
	if src == "" {
		return ""
	}
	out, err := format.Source([]byte(src))
	if err != nil {
		return src
	}
	return string(out)
}

// FormatGoTypeDecl applies gofmt to a single type declaration (e.g. "type Foo struct { ... }").
func FormatGoTypeDecl(decl string) string {
	return formatWrappedDecl(decl)
}

// formatWrappedDecl applies gofmt to a single declaration (type or func) fragment.
func formatWrappedDecl(decl string) string {
	decl = strings.TrimSpace(decl)
	if decl == "" {
		return ""
	}
	wrapped := "package x\n\n" + decl
	out, err := format.Source([]byte(wrapped))
	if err != nil {
		return decl
	}
	s := string(out)
	s = strings.TrimPrefix(s, "package x")
	return strings.TrimSpace(s)
}

// DedentCommonLeadingWhitespace removes the longest string of spaces and tabs that is a
// prefix of every non-blank line (textwrap.dedent-style). Blank lines are left as-is
// except they are still trimmed of that prefix when it matches.
func DedentCommonLeadingWhitespace(s string) string {
	lines := strings.Split(s, "\n")
	common := ""
	started := false
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		pref := leadingSpaceTabPrefix(line)
		if !started {
			common = pref
			started = true
			continue
		}
		common = sharedSpaceTabPrefix(common, pref)
		if common == "" {
			break
		}
	}
	if !started || common == "" {
		return strings.TrimRight(s, "\n\r")
	}
	var b strings.Builder
	for _, line := range lines {
		if strings.HasPrefix(line, common) {
			b.WriteString(strings.TrimPrefix(line, common))
		} else {
			b.WriteString(line)
		}
		b.WriteByte('\n')
	}
	return strings.TrimRight(b.String(), "\n\r")
}

func leadingSpaceTabPrefix(s string) string {
	i := 0
	for i < len(s) {
		switch s[i] {
		case ' ', '\t':
			i++
		default:
			return s[:i]
		}
	}
	return s[:i]
}

func sharedSpaceTabPrefix(a, b string) string {
	if len(a) > len(b) {
		a, b = b, a
	}
	i := 0
	for i < len(a) && i < len(b) && a[i] == b[i] {
		i++
	}
	return a[:i]
}
