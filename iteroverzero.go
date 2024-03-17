package sundrylint

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func LintIterOverZero(pass *analysis.Pass, node *ast.RangeStmt, stack []ast.Node) (ds []analysis.Diagnostic) {
	if _, ok := pass.TypesInfo.TypeOf(node.X).(*types.Slice); !ok {
		return nil
	}
	stmts, ncalls := restOfBlock(stack)
	if ncalls > 1 {
		return nil
	}
	if len(stmts) < 2 {
		return nil
	}
	/*
				if !isHTTPFuncOrMethodOnClient(pass.TypesInfo, call) {
					return true // the function call is not related to this check.
				}

				// Find the innermost containing block, and get the list
				// of statements starting with the one containing call.
				stmts, ncalls := restOfBlock(stack)
				if len(stmts) < 2 {
					// The call to the http function is the last statement of the block.
					return true
				}

				// Skip cases in which the call is wrapped by another (#52661).
				// Example:  resp, err := checkError(http.Get(url))
				if ncalls > 1 {
					return true
				}

				asg, ok := stmts[0].(*ast.AssignStmt)
				if !ok {
					return true // the first statement is not assignment.
				}

				resp := rootIdent(asg.Lhs[0])
				if resp == nil {
					return true // could not find the http.Response in the assignment.
				}

				def, ok := stmts[1].(*ast.DeferStmt)
				if !ok {
					return true // the following statement is not a defer.
				}
				root := rootIdent(def.Call.Fun)
				if root == nil {
					return true // could not find the receiver of the defer call.
				}

				if resp.Obj == root.Obj {
					pass.ReportRangef(root, "using %s before checking for errors", resp.Name)
				}
				return true
			})

		inspect.WithStack(nodeFilter, func(n ast.Node, push bool, stack []ast.Node) bool {
	*/

	return nil
}
