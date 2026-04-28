package parser

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"docgen/utils"
)

// OrderedSource is one non-test .go file in a category, parsed and ordered for
// docs, sidebar, and readme (see OrderedSources).
type OrderedSource struct {
	Path    string
	Meta    map[string]string
	DocType string
	Title   string
}

// DocType returns meta["type"] or "problem" when unset.
func DocType(meta map[string]string) string {
	if t := strings.TrimSpace(meta["type"]); t != "" {
		return t
	}
	return "problem"
}

// OrderKeyFromMeta returns the primary sort key: @index as integer when set and
// valid, otherwise leading digits in title (see utils.OrderPrefixFromTitle).
func OrderKeyFromMeta(meta map[string]string, title string) int {
	if idx := strings.TrimSpace(meta["index"]); idx != "" {
		if n, err := strconv.Atoi(idx); err == nil {
			return n
		}
	}
	return utils.OrderPrefixFromTitle(title)
}

// OrderedSources returns non-test .go sources in doc/readme/sidebar order:
// @index first (numeric: 1, 2, 11, …), else leading number in @problem/@title,
// then title A→Z. Skips files with no resolvable title (same as docgen).
func OrderedSources(files []string) []OrderedSource {
	var out []OrderedSource
	for _, file := range files {
		if utils.IsTestFile(file) {
			continue
		}
		meta := ParseMetadata(file)
		docType := DocType(meta)
		title := ResolveTitle(meta, docType)
		if title == "" {
			fmt.Fprintf(os.Stderr, "skip (add @problem: or @title:): %s\n", file)
			continue
		}
		out = append(out, OrderedSource{
			Path:    file,
			Meta:    meta,
			DocType: docType,
			Title:   title,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		ki := OrderKeyFromMeta(out[i].Meta, out[i].Title)
		kj := OrderKeyFromMeta(out[j].Meta, out[j].Title)
		if ki != kj {
			return ki < kj
		}
		return strings.ToLower(out[i].Title) < strings.ToLower(out[j].Title)
	})
	return out
}
