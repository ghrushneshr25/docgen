package parser

import (
	"os"
	"strings"
)

func ParseMetadata(file string) map[string]string {
	meta := make(map[string]string)

	data, _ := os.ReadFile(file)
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "/*") || strings.HasPrefix(line, "package ") {
			break
		}

		if strings.HasPrefix(line, "//") {
			line = strings.TrimPrefix(line, "//")
			line = strings.TrimSpace(line)

			if strings.HasPrefix(line, "@") {
				line = strings.TrimPrefix(line, "@") // ✅ REMOVE @

				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])   // problem
					value := strings.TrimSpace(parts[1]) // Towers of Hanoi
					meta[key] = value
				}
			}
		}
	}

	return meta
}

func ResolveTitle(meta map[string]string, t string) string {
	if t == "concept" {
		return meta["title"]
	}
	if meta["problem"] != "" {
		return meta["problem"]
	}
	return meta["title"]
}
