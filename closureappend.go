package sundrylint

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

func LintClosureAppend(pass *analysis.Pass, node *ast.CallExpr, stack []ast.Node) (ds []analysis.Diagnostic) {
	// Only check append function calls
	if ident, ok := node.Fun.(*ast.Ident); !ok || ident.Name != "append" {
		return nil
	}

	// Check if we are inside a closure (function literal)
	funcLit, inClosure := findEnclosingFuncLit(stack)
	if !inClosure {
		return nil
	}

	// Check if we are inside a loop
	// This indicates potential concurrent access issues
	loopStmt := findEnclosingLoopStmt(stack)
	if loopStmt == nil {
		return nil
	}

	// Check if the append operation accesses variables defined outside the closure
	// This indicates potential issues with shared state
	if accessesExternalVariables(pass, node, funcLit, stack) {
		return []analysis.Diagnostic{
			{
				Pos:      node.Pos(),
				End:      node.End(),
				Category: LinterName,
				Message:  SubLinterClosureAppendMessage,
			},
		}
	}

	return nil
}

// findEnclosingFuncLit finds if we are inside a function literal (closure)
func findEnclosingFuncLit(stack []ast.Node) (*ast.FuncLit, bool) {
	for i := len(stack) - 1; i >= 0; i-- {
		if funcLit, ok := stack[i].(*ast.FuncLit); ok {
			return funcLit, true
		}
	}
	return nil, false
}

// isLoopVariableCapture checks if the append operation captures loop variables
func isLoopVariableCapture(pass *analysis.Pass, callExpr *ast.CallExpr, funcLit *ast.FuncLit, stack []ast.Node) bool {
	// Check if we are inside a loop (for, range)
	loopStmt := findEnclosingLoopStmt(stack)
	if loopStmt == nil {
		return false
	}

	// Check if append arguments reference loop variables
	for _, arg := range callExpr.Args {
		if referencesLoopVariable(pass, arg, loopStmt, funcLit) {
			return true
		}
	}

	return false
}

// findEnclosingLoopStmt finds the nearest enclosing loop statement
func findEnclosingLoopStmt(stack []ast.Node) ast.Stmt {
	for i := len(stack) - 1; i >= 0; i-- {
		switch node := stack[i].(type) {
		case *ast.ForStmt, *ast.RangeStmt:
			return node.(ast.Stmt)
		}
	}
	return nil
}

// referencesLoopVariable checks if an expression references a loop variable
func referencesLoopVariable(pass *analysis.Pass, expr ast.Expr, loopStmt ast.Stmt, funcLit *ast.FuncLit) bool {
	var loopVars []*ast.Ident

	// Extract loop variables based on loop type
	switch stmt := loopStmt.(type) {
	case *ast.ForStmt:
		if stmt.Init != nil {
			if assign, ok := stmt.Init.(*ast.AssignStmt); ok {
				for _, lhs := range assign.Lhs {
					if ident, ok := lhs.(*ast.Ident); ok {
						loopVars = append(loopVars, ident)
					}
				}
			}
		}
	case *ast.RangeStmt:
		if stmt.Key != nil {
			if ident, ok := stmt.Key.(*ast.Ident); ok && ident.Name != "_" {
				loopVars = append(loopVars, ident)
			}
		}
		if stmt.Value != nil {
			if ident, ok := stmt.Value.(*ast.Ident); ok && ident.Name != "_" {
				loopVars = append(loopVars, ident)
			}
		}
	}

	// Check if the expression references any loop variable
	var found bool
	ast.Inspect(expr, func(n ast.Node) bool {
		if ident, ok := n.(*ast.Ident); ok {
			for _, loopVar := range loopVars {
				if pass.TypesInfo.ObjectOf(ident) == pass.TypesInfo.ObjectOf(loopVar) {
					found = true
					return false // Stop further inspection
				}
			}
		}
		return true
	})

	return found
}

// accessesExternalVariables checks if the append operation accesses variables defined outside the closure
func accessesExternalVariables(pass *analysis.Pass, callExpr *ast.CallExpr, funcLit *ast.FuncLit, stack []ast.Node) bool {
	// Only check the first argument (slice) for external variable access
	// This is the slice being appended to, which should not be shared in concurrent contexts
	if len(callExpr.Args) == 0 {
		return false
	}

	firstArg := callExpr.Args[0]
	return isVariableDefinedOutsideClosure(pass, firstArg, funcLit)
}

// isVariableDefinedOutsideClosure checks if a variable is defined outside the given closure
func isVariableDefinedOutsideClosure(pass *analysis.Pass, expr ast.Expr, funcLit *ast.FuncLit) bool {
	var found bool

	ast.Inspect(expr, func(n ast.Node) bool {
		if ident, ok := n.(*ast.Ident); ok && ident.Obj != nil {
			// Check if this identifier's declaration is outside the closure
			if isDeclaredOutsideClosure(pass, ident, funcLit) {
				found = true
				return false // Stop further inspection
			}
		}
		return true
	})

	return found
}

// isDeclaredOutsideClosure checks if an identifier is declared outside the given closure
func isDeclaredOutsideClosure(pass *analysis.Pass, ident *ast.Ident, funcLit *ast.FuncLit) bool {
	obj := ident.Obj
	if obj == nil {
		return false
	}

	// Get the declaration position
	declPos := obj.Pos()
	if declPos == token.NoPos {
		return false
	}

	// Check if declaration is outside the function literal
	funcLitStart := funcLit.Pos()
	funcLitEnd := funcLit.End()

	return declPos < funcLitStart || declPos > funcLitEnd
}