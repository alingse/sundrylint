package sundrylint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	LinterName = "sundrylint"
	LinterDesc = "This is a linter that combines multiple small checks, primarily derived from real-world errors encountered during the development process."

	SubLinterTimeparse        = `timeparse`
	SubLinterIteroverzero     = `iteroverzero`
	SubLinterFuncresultunused = `funcresultunused`
	SubLinterRangeappendall   = `rangeappendall`

	SubLinterRangeappendallMessage = `append all its data while range it`
	SubLinterAppendNoAssignMessage = `call strconv.AppendX but not keep func result`
	SubLinterMustCompileOutMessage = `call regexp.MustCompile with constant should be moved out of func`
)

type LinterSetting struct{}

func NewAnalyzer(setting LinterSetting) (*analysis.Analyzer, error) {
	a, err := newAnalyzer(setting)
	if err != nil {
		return nil, err
	}

	return &analysis.Analyzer{
		Name: LinterName,
		Doc:  LinterDesc,
		Run:  a.run,
		Requires: []*analysis.Analyzer{
			inspect.Analyzer,
			buildssa.Analyzer,
		},
	}, nil
}

type analyzer struct {
	setting LinterSetting
}

func newAnalyzer(setting LinterSetting) (*analyzer, error) {
	a := &analyzer{setting: setting}
	return a, nil
}

func (a *analyzer) run(pass *analysis.Pass) (interface{}, error) {
	_, _ = a.checkInspect(pass)
	return nil, nil
}

func (a *analyzer) checkInspect(pass *analysis.Pass) (interface{}, error) {
	inspectorInfo := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	checkNodes := []ast.Node{
		(*ast.CallExpr)(nil),
		(*ast.RangeStmt)(nil),
		(*ast.AssignStmt)(nil),
		(*ast.ReturnStmt)(nil),
		(*ast.DeferStmt)(nil),
		(*ast.KeyValueExpr)(nil),
		(*ast.CompositeLit)(nil),
	}
	inspectorInfo.WithStack(checkNodes, func(n ast.Node, push bool, stack []ast.Node) (proceed bool) {
		a.process(pass, n, push, stack)
		return true
	})
	return nil, nil
}

func (a *analyzer) process(pass *analysis.Pass, n ast.Node, push bool, stack []ast.Node) {
	if !push {
		return
	}

	switch node := n.(type) {
	case *ast.CallExpr:
		a.report(pass, LintTimeParseWrongArgsOrder(pass, node))
		//a.report(pass, LintFuncResultUnused(pass, node, stack))
		a.report(pass, LintRangeAppendAll(pass, node, stack))
		a.report(pass, AppendNoAssign(pass, node, stack))
		a.report(pass, MustCompileOut(pass, node, stack))
	case *ast.RangeStmt:
		a.report(pass, LintIterOverZero(pass, node, stack))
	}
}

func (a *analyzer) report(pass *analysis.Pass, ds []analysis.Diagnostic) {
	for _, d := range ds {
		pass.Report(d)
	}
}
