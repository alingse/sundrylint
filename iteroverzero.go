package sundrylint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

func LintIterOverZero(pass *analysis.Pass, node *ast.RangeStmt, stack []ast.Node) (ds []analysis.Diagnostic) {
	return nil
}
