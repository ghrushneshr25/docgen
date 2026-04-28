package utils

import "strings"

// NoOrderPrefix is returned by OrderPrefixFromTitle when the title does not
// start with digits (those entries sort after numbered titles, then by title).
const NoOrderPrefix = 1 << 20

// OrderPrefixFromTitle parses a leading decimal integer from the title for stable
// ordering (e.g. "1 Foo" → 1, "11 Bar" → 11, "01 Baz" → 1). Put order in @problem /
// @title as "1 …", "2 …", … so sidebar and category index sort numerically, not as strings.
func OrderPrefixFromTitle(title string) int {
	title = strings.TrimSpace(title)
	if title == "" {
		return NoOrderPrefix
	}
	n := 0
	i := 0
	for i < len(title) && i < 12 && title[i] >= '0' && title[i] <= '9' {
		n = n*10 + int(title[i]-'0')
		i++
	}
	if i == 0 {
		return NoOrderPrefix
	}
	return n
}
