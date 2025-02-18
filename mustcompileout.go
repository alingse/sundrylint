package sundrylint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

func MustCompileOut(pass *analysis.Pass, node *ast.CallExpr, stack []ast.Node) (ds []analysis.Diagnostic) {
	if !IsFuncPkg(pass, node, "regexp") {
		return nil
	}
	sel, ok := node.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "MustCompile" {
		return nil
	}
	if len(node.Args) != 1 {
		return nil
	}
	if !IsConst(pass, node.Args[0]) {
		return nil
	}

	if len(stack) > 1 {
		for _, parent := range stack[:len(stack)-1] {
			switch parentNode := parent.(type) {
			case *ast.FuncDecl:
				if parentNode.Name.Name == "init" {
					return nil
				}
				return []analysis.Diagnostic{
					{
						Pos:      node.Lparen,
						End:      node.Rparen,
						Category: LinterName,
						Message:  SubLinterMustCompileOutMessage,
					},
				}
			}
		}
	}
	return nil
}
