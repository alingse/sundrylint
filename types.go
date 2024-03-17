package sundrylint

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func IsVar(pass *analysis.Pass, e ast.Expr) bool {
	tv := pass.TypesInfo.Types[e]
	if tv.Addressable() && tv.Assignable() {
		return true
	}
	return false
}

func IsConst(pass *analysis.Pass, e ast.Expr) bool {
	switch e.(type) {
	case *ast.SelectorExpr, *ast.BasicLit, *ast.Ident:
	default:
		return false
	}

	tv := pass.TypesInfo.Types[e]
	if tv.Addressable() || tv.Assignable() {
		return false
	}
	return true
}

func IsPkg(pass *analysis.Pass, e ast.Expr, pkgPath string) bool {
	if selectExpr, ok := e.(*ast.SelectorExpr); ok {
		obj := pass.TypesInfo.ObjectOf(selectExpr.Sel)
		if obj == nil || obj.Pkg().Path() != pkgPath {
			return false
		}
		return true
	}
	return false
}

type FuncType struct {
	ArgsNum   int
	Signature string
}

func IsFunc(pass *analysis.Pass, node *ast.CallExpr, fnType FuncType) bool {
	if len(node.Args) != fnType.ArgsNum {
		return false
	}

	sign, ok := pass.TypesInfo.TypeOf(node.Fun).(*types.Signature)
	if !ok || sign.String() != fnType.Signature {
		return false
	}
	return true
}
