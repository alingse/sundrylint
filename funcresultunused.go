package sundrylint

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

const FuncResultUnusedMessage = "func result unused"

func LintFuncResultUnused(pass *analysis.Pass, node *ast.CallExpr, stack []ast.Node) (ds []analysis.Diagnostic) {
	// check sign
	sign, ok := pass.TypesInfo.TypeOf(node.Fun).(*types.Signature)
	if !ok || sign.Recv() != nil {
		return nil
	}
	if sign.Results().Len() == 0 {
		return nil
	}
	if !IsTupleAll(sign.Params(), IsBasicType) {
		return nil
	}

	// check is not a obj.Func
	if se, ok := node.Fun.(*ast.SelectorExpr); ok {
		fnObj := pass.TypesInfo.ObjectOf(se.Sel)
		// has a receiver
		if sign, ok := fnObj.Type().(*types.Signature); ok && sign.Recv() != nil {
			return nil
		}
		// not a call from pkg.XXX
		if xtype, ok := pass.TypesInfo.TypeOf(se.X).(*types.Basic); !ok || xtype.Kind() != types.Invalid {
			return nil
		}
	}

	if len(stack) < 2 {
		return nil
	}
	lastNode := stack[len(stack)-2]
	switch lastNode.(type) {
	case *ast.AssignStmt, *ast.KeyValueExpr, *ast.ReturnStmt:
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
	default:
	}
	return nil
}
