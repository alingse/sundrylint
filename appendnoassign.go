package sundrylint

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func AppendNoAssign(pass *analysis.Pass, node *ast.CallExpr, stack []ast.Node) (ds []analysis.Diagnostic) {
	if !IsFuncPkg(pass, node, "strconv") {
		return nil
	}
	sign, ok := pass.TypesInfo.TypeOf(node.Fun).(*types.Signature)
	if !ok {
		return nil
	}
	if sign.Results().Len() != 1 {
		return nil
	}
	if sign.Results().At(0).Type().String() != "[]byte" {
		return nil
	}
	if len(stack) > 1 {
		parent := stack[len(stack)-2]
		switch parentNode := parent.(type) {
		case *ast.AssignStmt, *ast.ReturnStmt:
			return nil
		case *ast.KeyValueExpr, *ast.CompositeLit:
			return nil
		case *ast.CallExpr:
			for _, arg := range parentNode.Args {
				if arg == node {
					return nil
				}
			}
		}
	}

	return []analysis.Diagnostic{
		{
			Pos:      node.Lparen,
			End:      node.Rparen,
			Category: LinterName,
			Message:  SubLinterAppendNoAssignMessage,
		},
	}
}
