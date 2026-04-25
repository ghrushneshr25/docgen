package parser

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// SubtestInfo is one t.Run(...) block from a test file, for documentation.
type SubtestInfo struct {
	Name string // first argument to t.Run
	Desc string // from // @desc: ... attached to the subtest body
	Code string // exact source of the t.Run(...) expression from the test file
}

// TestFilePath returns the path to test/<stem>_test.go next to a non-test .go source file.
func TestFilePath(sourceFile string) (string, bool) {
	dir := filepath.Dir(sourceFile)
	base := filepath.Base(sourceFile)
	if !strings.HasSuffix(base, ".go") || strings.HasSuffix(base, "_test.go") {
		return "", false
	}
	stem := strings.TrimSuffix(base, ".go")
	p := filepath.Join(dir, "test", stem+"_test.go")
	if _, err := os.Stat(p); err != nil {
		return "", false
	}
	return p, true
}

// ParseSubtests walks a *_test.go file and extracts each top-level t.Run inside Test* functions.
func ParseSubtests(testPath string) ([]SubtestInfo, error) {
	src, err := os.ReadFile(testPath)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, testPath, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	cm := ast.NewCommentMap(fset, file, file.Comments)

	var out []SubtestInfo
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil || !strings.HasPrefix(fn.Name.Name, "Test") {
			continue
		}
		for _, stmt := range fn.Body.List {
			es, ok := stmt.(*ast.ExprStmt)
			if !ok {
				continue
			}
			call, ok := es.X.(*ast.CallExpr)
			if !ok || !isTRunCall(call) {
				continue
			}
			name, err := stringLitValue(call.Args[0])
			if err != nil {
				continue
			}
			fl, ok := call.Args[1].(*ast.FuncLit)
			if !ok {
				continue
			}
			desc := descFromFuncLit(cm, fl)
			code := stripDescDirectiveLines(exprStmtSource(src, fset, es))
			code = DedentCommonLeadingWhitespace(code)
			if code == "" {
				continue
			}
			out = append(out, SubtestInfo{
				Name: name,
				Desc: desc,
				Code: code,
			})
		}
	}
	return out, nil
}

func isTRunCall(ce *ast.CallExpr) bool {
	if len(ce.Args) < 2 {
		return false
	}
	sel, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel == nil || sel.Sel.Name != "Run" {
		return false
	}
	// *testing.T is usually named "t"; allow any selector receiver (e.g. tt.Run in table tests).
	_, ok = sel.X.(*ast.Ident)
	return ok
}

func stringLitValue(e ast.Expr) (string, error) {
	bl, ok := e.(*ast.BasicLit)
	if !ok || bl.Kind != token.STRING {
		return "", errors.New("not a string literal")
	}
	return strconv.Unquote(bl.Value)
}

func descFromFuncLit(cm ast.CommentMap, fl *ast.FuncLit) string {
	if fl.Body == nil {
		return ""
	}
	for _, stmt := range fl.Body.List {
		for _, cg := range cm[stmt] {
			if d := parseDescFromCommentGroup(cg); d != "" {
				return d
			}
		}
	}
	return ""
}

// stripDescDirectiveLines removes whole-line // @desc: … comments from a source snippet
// (they are still shown from SubtestInfo.Desc above the code block).
func stripDescDirectiveLines(s string) string {
	lines := strings.Split(s, "\n")
	var b strings.Builder
	for _, line := range lines {
		if isDescDirectiveLine(line) {
			continue
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return strings.TrimRight(b.String(), "\n\r")
}

func isDescDirectiveLine(line string) bool {
	t := strings.TrimSpace(line)
	if !strings.HasPrefix(t, "//") {
		return false
	}
	rest := strings.TrimSpace(strings.TrimPrefix(t, "//"))
	return strings.HasPrefix(rest, "@desc:")
}

// exprStmtSource returns the source for the full t.Run(...) line(s). It begins at the
// start of the physical line (so leading tabs before t.Run are included); otherwise the
// first line would lose indentation while inner lines kept file-relative tabs.
func exprStmtSource(src []byte, fset *token.FileSet, stmt *ast.ExprStmt) string {
	start := fset.Position(stmt.Pos()).Offset
	end := fset.Position(stmt.End()).Offset
	if start < 0 || end > len(src) || start > end {
		return ""
	}
	from := lineByteStart(src, start)
	for i := from; i < start; i++ {
		if src[i] != ' ' && src[i] != '\t' {
			from = start
			break
		}
	}
	return strings.TrimRight(string(src[from:end]), "\n\r")
}

func lineByteStart(src []byte, pos int) int {
	if pos <= 0 {
		return 0
	}
	if pos > len(src) {
		pos = len(src)
	}
	for i := pos - 1; i >= 0; i-- {
		if src[i] == '\n' {
			return i + 1
		}
	}
	return 0
}

func parseDescFromCommentGroup(cg *ast.CommentGroup) string {
	for _, c := range cg.List {
		text := strings.TrimSpace(c.Text)
		text = strings.TrimPrefix(text, "//")
		text = strings.TrimSpace(text)
		if strings.HasPrefix(text, "@desc:") {
			return strings.TrimSpace(strings.TrimPrefix(text, "@desc:"))
		}
	}
	return ""
}
