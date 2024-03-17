package sundrylint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	LinterName = "sundrylint"
	LinterDesc = "This is a linter that combines multiple small checks, primarily derived from real-world errors encountered during the development process."
)

type LinterSetting struct{}

func NewAnalyzer(setting LinterSetting) (*analysis.Analyzer, error) {
	a, err := newAnalyzer(setting)
	if err != nil {
		return nil, err
	}

	return &analysis.Analyzer{
		Name:     LinterName,
		Doc:      LinterDesc,
		Run:      a.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}, nil
}

type analyzer struct {
	setting LinterSetting
}

func newAnalyzer(setting LinterSetting) (*analyzer, error) {
	a := &analyzer{
		setting: setting,
	}
	return a, nil
}

func (a *analyzer) run(pass *analysis.Pass) (interface{}, error) {
	inspectorInfo := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	checkNodes := []ast.Node{
		(*ast.CallExpr)(nil),
		(*ast.RangeStmt)(nil),
	}
	inspectorInfo.WithStack(checkNodes, func(n ast.Node, push bool, stack []ast.Node) (proceed bool) {
		a.process(pass, n, push, stack)
		return true
	})
	return nil, nil
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

func (a *analyzer) process(pass *analysis.Pass, n ast.Node, push bool, stack []ast.Node) {
	if !push {
		return
	}

	switch node := n.(type) {
	case *ast.CallExpr:
		a.report(pass, LintTimeParseWrongArgsOrder(pass, node))
	case *ast.RangeStmt:
		a.report(pass, LintIterOverZero(pass, node, stack))
	}
}

func (a *analyzer) report(pass *analysis.Pass, ds []analysis.Diagnostic) {
	for _, d := range ds {
		pass.Report(d)
	}
}
