package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"regexp"
	"strings"
)

// Banner blocks look like the middle line of
//
//	// ================================
//	// Traversal
//	// ================================
var anyWhitespace = regexp.MustCompile(`\s+`)
var reOperationTag = regexp.MustCompile(`^\s*//\s*@operation\s+(.+)$`)
var reStructureTag = regexp.MustCompile(`^\s*//\s*@structure\s+(.+)$`)

// NormalizeOpTitle maps doc subsection / banner titles to a single key for matching.
func NormalizeOpTitle(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, "-", " ")
	return anyWhitespace.ReplaceAllString(s, " ")
}

// titleFromBannerSubline returns the title from a line like "//  Traversal  ".
func titleFromBannerSubline(line string) (string, bool) {
	t := strings.TrimSpace(line)
	if !strings.HasPrefix(t, "//") {
		return "", false
	}
	rest := strings.TrimSpace(t[2:])
	if rest == "" || rest == "================================" {
		return "", false
	}
	return rest, true
}

// isOperationTag reports whether line is a // @operation Name directive (docgen).
func isOperationTag(line string) bool {
	_, ok := titleFromOperationTagLine(line)
	return ok
}

// titleFromOperationTagLine returns the operation title from a line like "// @operation Traversal".
func titleFromOperationTagLine(line string) (string, bool) {
	m := reOperationTag.FindStringSubmatch(line)
	if len(m) < 2 {
		return "", false
	}
	return strings.TrimSpace(m[1]), true
}

// isStructureTag reports whether line is a // @structure Name directive (mark before a type, docgen).
func isStructureTag(line string) bool {
	_, ok := titleFromStructureTagLine(line)
	return ok
}

// titleFromStructureTagLine returns the type name from a line like "// @structure ListNode".
func titleFromStructureTagLine(line string) (string, bool) {
	m := reStructureTag.FindStringSubmatch(line)
	if len(m) < 2 {
		return "", false
	}
	return strings.TrimSpace(m[1]), true
}

// goSectionsFromOperationDirectives returns line ranges for each // @operation Title line
// in file order. Content is from the line after the tag through the line before the next
// @operation, @operation tag line, or EOF.
func goSectionsFromOperationDirectives(lines []string) (sections []struct{ title, key string; start, end int }) {
	for i, line := range lines {
		title, ok := titleFromOperationTagLine(line)
		if !ok {
			continue
		}
		start0 := i + 1
		if start0 >= len(lines) {
			continue
		}
		end0 := len(lines) - 1
		for j := i + 1; j < len(lines); j++ {
			if isOperationTag(lines[j]) {
				end0 = j - 1
				break
			}
			if isStructureTag(lines[j]) {
				end0 = j - 1
				break
			}
		}
		if end0 < start0 {
			end0 = start0
		}
		sections = append(sections, struct{ title, key string; start, end int }{
			title: title,
			key:   NormalizeOpTitle(title),
			start: start0,
			end:   end0,
		})
	}
	return sections
}

// goSectionsFromBanners returns line ranges and titles for "banner" section blocks
// in the .go file (e.g. before methods). Used to place functions into subsections.
func goSectionsFromBanners(lines []string) (sections []struct{ title, key string; start, end int }) {
	const banner = "// ================================"

	for i := 0; i+2 < len(lines); i++ {
		if strings.TrimSpace(lines[i]) != banner {
			continue
		}
		title, ok1 := titleFromBannerSubline(lines[i+1])
		if !ok1 || strings.TrimSpace(lines[i+2]) != banner {
			continue
		}
		// code / comments start at i+3; section ends on line before the next top-of-file banner
		var end int
		if i+3 >= len(lines) {
			end = len(lines) - 1
		} else {
			for j := i + 3; j < len(lines); j++ {
				if strings.TrimSpace(lines[j]) == banner {
					end = j - 1
					break
				}
				if isOperationTag(lines[j]) {
					end = j - 1
					break
				}
				if isStructureTag(lines[j]) {
					end = j - 1
					break
				}
				end = j
			}
		}
		if end < i+3 {
			end = i + 3
		}
		sections = append(sections, struct{ title, key string; start, end int }{
			title: title,
			key:   NormalizeOpTitle(title),
			start: i + 3,
			end:   end,
		})
	}
	return sections
}

// OperationDirectiveCount returns how many // @operation directives appear in the file (in order).
func OperationDirectiveCount(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	lines := strings.Split(string(data), "\n")
	return len(goSectionsFromOperationDirectives(lines)), nil
}

// OperationSnippet is one gofmt'd function under // @operation, in source order.
type OperationSnippet struct {
	Title  string
	Key    string
	Source string
}

// ExtractOperationCodeOrdered returns function source for each // @operation in file order.
func ExtractOperationCodeOrdered(path string) ([]OperationSnippet, error) {
	byKey, err := ExtractOperationCodeBySection(path)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	ranges := goSectionsFromOperationDirectives(lines)
	out := make([]OperationSnippet, 0, len(ranges))
	for _, r := range ranges {
		src := strings.TrimSpace(byKey[r.key])
		out = append(out, OperationSnippet{Title: r.title, Key: r.key, Source: src})
	}
	return out, nil
}

func sectionKeyForLine(sections []struct{ title, key string; start, end int }, line1Based int) string {
	line0 := line1Based - 1
	var key string
	for _, s := range sections {
		if line0 >= s.start && line0 <= s.end {
			// if overlapping (shouldn't), prefer last matching
			if s.key != "" {
				key = s.key
			}
		}
	}
	return key
}

// ExtractOperationCodeBySection returns normalized subsection key -> gofmt'd source for
// that section, from // @operation <Title> (and legacy // === // Title // ===).
// // @structure <Name> is treated as a boundary (same as // @operation) for banner spans.
// Keys match NormalizeOpTitle and @subsection: in the leading doc block.
func ExtractOperationCodeBySection(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	// Banner sections (// ===) first, then // @operation ranges so overlapping lines take the
	// @operation key (e.g. methods after a struct banner block that ends at @operation).
	sectionRanges := goSectionsFromBanners(lines)
	sectionRanges = append(sectionRanges, goSectionsFromOperationDirectives(lines)...)

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, data, 0)
	if err != nil {
		return nil, err
	}

	byKey := make(map[string][]string)

	for _, decl := range node.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name == nil || fn.Name.Name == "init" {
			continue
		}
		l := fset.Position(fn.Pos()).Line
		key := sectionKeyForLine(sectionRanges, l)
		if key == "" {
			continue
		}
		var buf bytes.Buffer
		if err := printer.Fprint(&buf, fset, fn); err != nil {
			return nil, fmt.Errorf("print func %q: %w", fn.Name.Name, err)
		}
		s := strings.TrimSpace(buf.String())
		s = formatWrappedDecl(s)
		if s == "" {
			continue
		}
		byKey[key] = append(byKey[key], s)
	}
	out := make(map[string]string)
	for k, parts := range byKey {
		out[k] = strings.Join(parts, "\n\n")
	}
	return out, nil
}
