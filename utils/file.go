package utils

import "strings"

func IsTestFile(file string) bool {
	return strings.HasSuffix(file, "_test.go")
}
