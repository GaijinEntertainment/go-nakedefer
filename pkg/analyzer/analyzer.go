package analyzer

import (
	"bytes"
	"errors"
	"flag"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var (
	ErrEmptyPattern = errors.New("pattern can't be empty")
)

type analyzer struct {
	typesInfo *types.Info
	exclude   PatternsList
}

// NewAnalyzer returns a go/analysis-compatible analyzer.
func NewAnalyzer(exclude []string) (*analysis.Analyzer, error) {
	a := analyzer{} //nolint:exhaustruct

	var err error

	a.exclude, err = newPatternsList(exclude)
	if err != nil {
		return nil, err
	}

	return &analysis.Analyzer{ //nolint:exhaustruct
		Name:     "nackedefer",
		Doc:      "Checks that deferred call does not return anything.",
		Run:      a.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Flags:    a.newFlagSet(),
	}, nil
}

func (a *analyzer) newFlagSet() flag.FlagSet {
	fs := flag.NewFlagSet("nackedefer flags", flag.PanicOnError)

	fs.Var(
		&reListVar{values: &a.exclude},
		"e",
		"Regular expression to exclude function names",
	)

	return *fs
}

func (a *analyzer) run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector) //nolint:forcetypeassert

	nodeFilter := []ast.Node{
		(*ast.DeferStmt)(nil),
	}

	a.typesInfo = pass.TypesInfo

	insp.Preorder(nodeFilter, a.newVisitor(pass))

	return nil, nil //nolint:nilnil
}

func (a *analyzer) newVisitor(pass *analysis.Pass) func(node ast.Node) {
	return func(node ast.Node) {
		deferStmt, ok := node.(*ast.DeferStmt)
		if !ok {
			return
		}

		if deferStmt.Call == nil {
			return
		}

		funcName := a.funcName(deferStmt.Call)
		if funcName != "" && a.exclude.MatchesAny(funcName) {
			return
		}

		var isFuncReturnType bool
		switch v := deferStmt.Call.Fun.(type) {
		case *ast.FuncLit: // function is anonymous
			isFuncReturnType = a.isFuncLitReturnType(v)
		case *ast.Ident:
			isFuncReturnType = a.isIdentReturnType(v)
		case *ast.SelectorExpr:
			isFuncReturnType = a.isSelExprReturnType(v)
		default:
			return
		}

		if !isFuncReturnType {
			return
		}

		pass.Reportf(node.Pos(), "deferred call should not return anything")
	}
}

func (a *analyzer) isIdentReturnType(ident *ast.Ident) bool {
	if ident == nil || ident.Obj == nil {
		return false
	}

	funcDecl, ok := ident.Obj.Decl.(*ast.FuncDecl)
	if !ok {
		return false
	}

	if funcDecl.Type == nil || funcDecl.Type.Results == nil {
		return false
	}

	if len(funcDecl.Type.Results.List) == 0 {
		return false
	}

	return true
}

func (a *analyzer) isFuncLitReturnType(funcLit *ast.FuncLit) bool {
	if funcLit == nil || funcLit.Type == nil {
		return false
	}

	if funcLit.Type == nil || funcLit.Type.Results == nil {
		return false
	}

	if len(funcLit.Type.Results.List) == 0 {
		return false
	}

	return true
}

func (a *analyzer) isSelExprReturnType(selExpr *ast.SelectorExpr) bool {
	if selExpr == nil {
		return false
	}

	t, ok := a.typesInfo.Types[selExpr].Type.(*types.Signature)
	if !ok {
		return false
	}

	if t.Results() == nil || t.Results().Len() == 0 {
		return false
	}

	return true
}

func (a *analyzer) funcName(call *ast.CallExpr) string {
	fn, ok := a.getFunc(call)
	if !ok {
		return gofmt(call.Fun)
	}

	name := fn.FullName()
	name = strings.ReplaceAll(name, "(", "")
	name = strings.ReplaceAll(name, ")", "")

	return name
}

func (a *analyzer) getFunc(call *ast.CallExpr) (*types.Func, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}

	fn, ok := a.typesInfo.ObjectOf(sel.Sel).(*types.Func)
	if !ok {
		return nil, false
	}

	return fn, true
}

func gofmt(x interface{}) string {
	buf := bytes.Buffer{}
	fs := token.NewFileSet()
	printer.Fprint(&buf, fs, x)

	return buf.String()
}
