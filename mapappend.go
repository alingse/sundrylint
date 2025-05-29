package sundrylint

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func LintMapAppend(pass *analysis.Pass, node *ast.CallExpr, stack []ast.Node) (ds []analysis.Diagnostic) {
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

	// Only check append function calls
	if funcName != "append" {
		return
	}

	// Check if in assignment statement
	assignStmt, ok := findParentAssignStmt(node, stack)
	if !ok {
		return
	}

	// Check if left side is map index expression
	lhsMap, lhsKey, ok := isMapIndexExpr(assignStmt.Lhs[0])
	if !ok {
		return
	}

	// Check if first argument is map index expression
	rhsMap, rhsKey, ok := isMapIndexExpr(node.Args[0])
	if !ok {
		return
	}

	// Check if both sides reference same map
	lhsMapCode, err := GetCode(pass.Fset, lhsMap)
	if err != nil {
		return
	}
	rhsMapCode, err := GetCode(pass.Fset, rhsMap)
	if err != nil {
		return
	}

	if lhsMapCode != rhsMapCode {
		return
	}

	// Check if keys are different
	lhsKeyCode, err := GetCode(pass.Fset, lhsKey)
	if err != nil {
		return
	}
	rhsKeyCode, err := GetCode(pass.Fset, rhsKey)
	if err != nil {
		return
	}

	if lhsKeyCode == rhsKeyCode {
		return
	}

	// Generate diagnostic message
	return append(ds, analysis.Diagnostic{
		Pos:      node.Pos(),
		End:      node.End(),
		Category: LinterName,
		Message:  "potential incorrect map append: using different keys for append and assignment",
	})
}

// Find parent assignment statement
func findParentAssignStmt(f ast.Node, stack []ast.Node) (*ast.AssignStmt, bool) {
	for i := len(stack) - 1; i >= 0; i-- {
		if assign, ok := stack[i].(*ast.AssignStmt); ok {
			return assign, true
		}
	}
	return nil, false
}

// Check if expression is map index access
func isMapIndexExpr(expr ast.Expr) (mapExpr ast.Expr, keyExpr ast.Expr, ok bool) {
	indexExpr, ok := expr.(*ast.IndexExpr)
	if !ok {
		return nil, nil, false
	}
	return indexExpr.X, indexExpr.Index, true
}
