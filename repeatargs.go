package sundrylint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

func LintRepeatArgs(pass *analysis.Pass, node *ast.CallExpr) (ds []analysis.Diagnostic) {
	argsMap := make(map[string][]any)
	for _, arg := range node.Args {
		_, ok := arg.(*ast.CallExpr)
		if !ok {
			continue
		}
		code, err := GetCode(pass.Fset, arg)
		if err != nil {
			continue
		}
		argsMap[code] = append(argsMap[code], arg)
	}
	for _, args := range argsMap {
		if len(args) > 1 {
			ds = append(ds, analysis.Diagnostic{
				Pos:      node.Pos(),
				End:      node.End(),
				Category: LinterName,
				Message:  SubLinterRepeatArgsMessage,
			})
		}
	}
	return ds
}
