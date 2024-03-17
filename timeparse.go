package sundrylint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const TimeParseLintMessage = "call func time.Parse may have incorrect parameters, potentially swapping the layout and value arguments."

const (
	timePkgPath = `time`
)

var timeParseFunc = FuncType{
	ArgsNum:   2,
	Signature: `func(layout string, value string) (time.Time, error)`,
}

func isTimeParseFunc(pass *analysis.Pass, node *ast.CallExpr) bool {
	return IsFunc(pass, node, timeParseFunc) && IsPkg(pass, node.Fun, timePkgPath)
}

func processTimeParse(pass *analysis.Pass, node *ast.CallExpr) (ds []analysis.Diagnostic) {
	if !isTimeParseFunc(pass, node) {
		return nil
	}
	if IsVar(pass, node.Args[0]) && IsConst(pass, node.Args[1]) {
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
