package sundrylint

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var noCheckFuncs = []string{
	"print",
}

func LintRepeatArgs(pass *analysis.Pass, node *ast.CallExpr) (ds []analysis.Diagnostic) {
	funcName, err := GetCode(pass.Fset, node.Fun)
	if err != nil {
		return nil
	}

	funcName = strings.ToLower(funcName)
	for _, fn := range noCheckFuncs {
		if strings.Contains(funcName, fn) {
			return
		}
	}

	argsMap := make(map[string][]any)
	for _, arg := range node.Args {
		argCall, ok := arg.(*ast.CallExpr)
		if !ok {
			continue
		}
		if len(argCall.Args) == 0 {
			continue
		}
		if _, ok := IsBuiltinFunc2(pass, argCall); ok {
			continue
		}
		if _, ok := IsTypeNameFunc(pass, argCall); ok {
			continue
		}
		allConst := true
		for _, arg := range argCall.Args {
			if !IsConst(pass, arg) {
				allConst = false
				break
			}
		}
		if allConst {
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
