package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// NewAnalyzer returns a go/analysis-compatible analyzer.
func NewAnalyzer() (*analysis.Analyzer, error) {
	return &analysis.Analyzer{ //nolint:exhaustruct
		Name:     "defer",
		Doc:      "Checks that defer statement defers a function that does not return any type.",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}, nil
}

func run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector) //nolint:forcetypeassert

	nodeFilter := []ast.Node{
		(*ast.DeferStmt)(nil),
	}

	insp.Preorder(nodeFilter, newVisitor(pass))

	return nil, nil //nolint:nilnil
}

func newVisitor(pass *analysis.Pass) func(node ast.Node) {
	return func(node ast.Node) {
		deferStmt, ok := node.(*ast.DeferStmt)
		if !ok {
			return
		}

		if deferStmt.Call == nil {
			return
		}

		var outgoingFieldList *ast.FieldList

		switch v := deferStmt.Call.Fun.(type) {
		case *ast.Ident: // function is named
			outgoingFieldList = getFuncDeclResults(v)
		case *ast.FuncLit: // function is anonymous
			outgoingFieldList = getFuncLitResults(v)
		default:
			return
		}

		if outgoingFieldList == nil || outgoingFieldList.List == nil {
			return
		}

		if len(outgoingFieldList.List) == 0 {
			return
		}

		pass.Reportf(node.Pos(), "deferred call should not return any type")
	}
}

func getFuncDeclResults(ident *ast.Ident) *ast.FieldList {
	if ident.Obj == nil {
		return nil
	}

	funcDecl, ok := ident.Obj.Decl.(*ast.FuncDecl)
	if !ok {
		return nil
	}

	return funcDecl.Type.Results
}

func getFuncLitResults(funcLit *ast.FuncLit) *ast.FieldList {
	if funcLit.Type == nil {
		return nil
	}

	return funcLit.Type.Results
}
