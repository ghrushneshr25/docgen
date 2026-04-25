package parser

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"
)

// ExtractFunctions returns formatted Go source for each top-level func (and method)
// in the file, excluding init.
func ExtractFunctions(path string) ([]string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, decl := range node.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name == nil || fn.Name.Name == "init" {
			continue
		}
		var buf bytes.Buffer
		if err := printer.Fprint(&buf, fset, fn); err != nil {
			continue
		}
		s := strings.TrimSpace(buf.String())
		s = formatWrappedDecl(s)
		if s != "" {
			out = append(out, s)
		}
	}
	return out, nil
}
