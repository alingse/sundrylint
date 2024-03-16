package sundrylint

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

const (
	timePkgPath = `time`

	timeParseFuncName = `Parse`
	timeParseFuncSign = `func(layout string, value string) (time.Time, error)`
	timeParseArgsNum  = 2

	TimeParseLintMessage = "call func time.Parse may have incorrect parameters, potentially swapping the layout and value arguments."
)

func isTimeParse(pass *analysis.Pass, node *ast.CallExpr) bool {
	if len(node.Args) != timeParseArgsNum {
		return false
	}

	fnType, ok := pass.TypesInfo.TypeOf(node.Fun).(*types.Signature)
	if !ok || fnType.String() != timeParseFuncSign {
		return false
	}
	// time.Parse()
	if selectExpr, ok := node.Fun.(*ast.SelectorExpr); ok {
		obj := pass.TypesInfo.ObjectOf(selectExpr.Sel)
		if obj == nil || obj.Pkg().Path() != timePkgPath || obj.Name() != timeParseFuncName {
			return false
		}
		return true
	}
	return false
}

func isVar(pass *analysis.Pass, e ast.Expr) bool {
	tv := pass.TypesInfo.Types[e]
	if tv.Addressable() && tv.Assignable() {
		return true
	}
	return false
}

func isTimeConst(pass *analysis.Pass, e ast.Expr) bool {
	tv := pass.TypesInfo.Types[e]
	if tv.Addressable() || tv.Assignable() {
		return false
	}
	// time.DateTime or others
	if selectExpr, ok := e.(*ast.SelectorExpr); ok {
		obj := pass.TypesInfo.ObjectOf(selectExpr.Sel)
		if obj == nil || obj.Pkg().Path() != timePkgPath {
			return false
		}
		return true
	}
	return false
}

func processTimeParse(pass *analysis.Pass, node *ast.CallExpr) (ds []analysis.Diagnostic) {
	isParse := isTimeParse(pass, node)
	if !isParse {
		return nil
	}
	if isVar(pass, node.Args[0]) && isTimeConst(pass, node.Args[1]) {
		return []analysis.Diagnostic{
			{
				Pos:      node.Pos(),
				End:      node.End(),
				Category: LinterName,
				Message:  TimeParseLintMessage,
			},
		}
	}
	return nil
}
