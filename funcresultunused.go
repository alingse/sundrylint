package sundrylint

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

const FuncResultUnusedMessage = "func result unused"

func LintFuncResultUnused(pass *analysis.Pass, node *ast.CallExpr, stack []ast.Node) (ds []analysis.Diagnostic) {
	if len(node.Args) == 0 {
		return nil
	}
	sign, ok := pass.TypesInfo.TypeOf(node.Fun).(*types.Signature)
	if !ok {
		return nil
	}
	if sign.Recv() != nil {
		return nil
	}

	if sign.Results().Len() == 0 || sign.Params().Len() == 0 {
		return nil
	}
	if !IsTupleAll(sign.Params(), IsBasicType) {
		return nil
	}
	if !IsTupleAll(sign.Results(), IsBasicType) {
		return nil
	}
	if len(stack) < 2 {
		return nil
	}
	lastNode := stack[len(stack)-2]
	switch lastNode.(type) {
	case *ast.AssignStmt:
		return nil
	case *ast.ExprStmt:
		return []analysis.Diagnostic{
			{
				Pos:      node.Pos(),
				End:      node.End(),
				Category: LinterName,
				Message:  FuncResultUnusedMessage,
			},
		}
	}
	return nil
}
