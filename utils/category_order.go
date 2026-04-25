package utils

import (
	"io/fs"
	"os"
	"sort"
	"strings"
)

// OrderedDirEntries returns directory entries from entries: names listed in orderFile
// appear first (in that order), then any remaining directories sorted by name.
// Missing lines in orderFile are ignored; unknown dirs not in the file are appended sorted.
// Empty or missing orderFile falls back to sorting all dirs by name.
func OrderedDirEntries(entries []fs.DirEntry, orderFile string) []fs.DirEntry {
	var dirs []fs.DirEntry
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e)
		}
	}
	if len(dirs) == 0 {
		return nil
	}
	pref := loadCategoryOrder(orderFile)
	if len(pref) == 0 {
		sort.Slice(dirs, func(i, j int) bool {
			return dirs[i].Name() < dirs[j].Name()
		})
		return dirs
	}

	byName := make(map[string]fs.DirEntry, len(dirs))
	for _, e := range dirs {
		byName[e.Name()] = e
	}

	var out []fs.DirEntry
	seen := make(map[string]bool)
	for _, name := range pref {
		name = strings.TrimSpace(name)
		if name == "" || strings.HasPrefix(name, "#") {
			continue
		}
		if e, ok := byName[name]; ok {
			out = append(out, e)
			seen[name] = true
		}
	}
	var rest []string
	for name := range byName {
		if !seen[name] {
			rest = append(rest, name)
		}
	}
	sort.Strings(rest)
	for _, name := range rest {
		out = append(out, byName[name])
	}
	return out
}

func loadCategoryOrder(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var lines []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lines = append(lines, line)
	}
	return lines
}
