package sundrylint

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func LintComparePtr(pass *analysis.Pass, node *ast.BinaryExpr) (ds []analysis.Diagnostic) {
	if node.Op != token.EQL && node.Op != token.NEQ {
		return nil
	}

	xType := pass.TypesInfo.TypeOf(node.X)
	yType := pass.TypesInfo.TypeOf(node.Y)

	if isPointerType(xType) && isPointerType(yType) {
		return []analysis.Diagnostic{
			{
				Pos:      node.Pos(),
				End:      node.End(),
				Category: LinterName,
				Message:  SubLinterComparePtrMessage,
			},
		}
	}

	return nil
}

func isPointerType(t types.Type) bool {
	_, ok := t.(*types.Pointer)
	return ok
}
