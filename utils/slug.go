package utils

import (
	"regexp"
	"strings"
)

// leadingOrderPrefix strips a leading order index from titles for URL slugs only:
// one or more digits, then spaces/dots/dashes (e.g. "1 ", "11 -", "01 ").
// Page title in MDX still shows the full @problem / @title string.
var leadingOrderPrefix = regexp.MustCompile(`^\d{1,12}[\s.-]+`)

func StripLeadingOrderPrefix(s string) string {
	s = strings.TrimSpace(s)
	return leadingOrderPrefix.ReplaceAllString(s, "")
}

func Slugify(s string) string {
	s = strings.ToLower(s)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	return strings.Trim(re.ReplaceAllString(s, "-"), "-")
}

// DocSlug is Slugify(StripLeadingOrderPrefix(title)) for filenames and Docusaurus doc ids.
func DocSlug(title string) string {
	return Slugify(StripLeadingOrderPrefix(title))
}

func FormatTitle(s string) string {
	s = strings.ReplaceAll(s, "-", " ")
	words := strings.Fields(s)
	for i, w := range words {
		if w == "" {
			continue
		}
		words[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
	}
	return strings.Join(words, " ")
}