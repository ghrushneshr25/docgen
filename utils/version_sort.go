package utils

import (
	"path/filepath"
	"sort"
	"strconv"
	"unicode"
)

// VersionLess compares two strings in "version" (natural) order so that
// problem2.go sorts before problem10.go.
func VersionLess(a, b string) bool {
	la, lb := len(a), len(b)
	ia, ib := 0, 0
	for ia < la && ib < lb {
		ca, cb := a[ia], b[ib]
		if isDigitByte(ca) && isDigitByte(cb) {
			ja, jb := ia, ib
			for ja < la && isDigitByte(a[ja]) {
				ja++
			}
			for jb < lb && isDigitByte(b[jb]) {
				jb++
			}
			na, errA := strconv.ParseUint(a[ia:ja], 10, 64)
			nb, errB := strconv.ParseUint(b[ib:jb], 10, 64)
			if errA != nil || errB != nil {
				return a[ia:ja] < b[ib:jb]
			}
			if na != nb {
				return na < nb
			}
			if (ja - ia) != (jb - ib) {
				return (ja - ia) < (jb - ib)
			}
			ia, ib = ja, jb
			continue
		}
		ra := unicode.ToLower(rune(ca))
		rb := unicode.ToLower(rune(cb))
		if ra != rb {
			return ra < rb
		}
		ia++
		ib++
	}
	return la < lb
}

func isDigitByte(c byte) bool {
	return c >= '0' && c <= '9'
}

// SortSourcePathsByBase sorts file paths by filepath.Base using VersionLess.
func SortSourcePathsByBase(paths []string) {
	sort.Slice(paths, func(i, j int) bool {
		return VersionLess(filepath.Base(paths[i]), filepath.Base(paths[j]))
	})
}
