package sundrylint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const TimeParseLintMessage = "call func time.Parse may have incorrect args order, potentially swapping the layout and value arguments."

const (
	timePkgPath = `time`
)

var timeParseFunc = FuncType{
	ArgsNum:    2,
	Signature:  `func(layout string, value string) (time.Time, error)`,
	ResultsNum: 2,
}

func isTimeParseFunc(pass *analysis.Pass, node *ast.CallExpr) bool {
	return IsFuncType(pass, node, timeParseFunc) && IsFuncPkg(pass, node, timePkgPath)
}

func LintTimeParseWrongArgsOrder(pass *analysis.Pass, node *ast.CallExpr) (ds []analysis.Diagnostic) {
	if !isTimeParseFunc(pass, node) {
		return nil
	}

	if IsVar(pass, node.Args[0]) && IsConst(pass, node.Args[1]) {
		// skip for case `time.Parse(layout, "2015")` in _test.go
		if IsTestFile(pass, node) && IsBasicLit(node.Args[1]) {
			return nil
		}

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
