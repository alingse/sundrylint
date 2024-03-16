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
	inspectorInfo.Preorder(nil, a.AsCheckVisitor(pass))
	return nil, nil
}

func (a *analyzer) AsCheckVisitor(pass *analysis.Pass) func(ast.Node) {
	return func(n ast.Node) {
		switch node := n.(type) {
		case *ast.CallExpr:
			ds := processTimeParse(pass, node)
			for _, d := range ds {
				pass.Report(d)
			}
		}
	}
}
