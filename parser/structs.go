package parser

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
)

type StructInfo struct {
	Name string
	Code string
	Doc  string // prose from @subsection matching Name (concepts)
}

func ExtractStructs(path string) ([]StructInfo, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}

	var result []StructInfo

	for _, decl := range node.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}

		for _, spec := range gen.Specs {
			ts := spec.(*ast.TypeSpec)
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}

			var buf bytes.Buffer
			printer.Fprint(&buf, fset, st)

			code := "type " + ts.Name.Name + " " + buf.String()
			code = FormatGoTypeDecl(code)

			result = append(result, StructInfo{
				Name: ts.Name.Name,
				Code: code,
			})
		}
	}

	return result, nil
}
