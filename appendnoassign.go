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
		switch parent.(type) {
		case *ast.AssignStmt:
			return nil
		default:
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
	//panic(stack[len(stack)-1])
	/*
		ast.Parent(node)
		panic(fmt.Sprintf("node %#v", node))
		/*
				  if len(parentStack) > 0 {
			            // 查看直接的父节点

			                // 若为赋值语句，说明函数调用的结果被赋值
			                fmt.Printf("CallExpr at %v is part of an assignment.\n", callExpr.Pos())
			            default:
			                // 其他情况：可能没有被赋值
			                fmt.Printf("CallExpr at %v is not directly assigned.\n", callExpr.Pos())
			            }
			        }
	*/
	return nil
}
