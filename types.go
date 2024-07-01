package sundrylint

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func IsVar(pass *analysis.Pass, e ast.Expr) bool {
	tv := pass.TypesInfo.Types[e]
	if tv.Addressable() && tv.Assignable() {
		return true
	}
	return false
}

func IsConst(pass *analysis.Pass, e ast.Expr) bool {
	switch e.(type) {
	case *ast.SelectorExpr, *ast.BasicLit, *ast.Ident:
	case *ast.TypeAssertExpr, *ast.CallExpr, *ast.IndexExpr, *ast.StarExpr:
		return false
	default:
		return false
	}

	tv := pass.TypesInfo.Types[e]
	if tv.Addressable() || tv.Assignable() {
		return false
	}
	return true
}

func GetFuncExprIdent(pass *analysis.Pass, e ast.Expr) *ast.Ident {
	switch et := e.(type) {
	case *ast.SelectorExpr:
		return et.Sel
	case *ast.Ident:
		return et
	}
	return nil
}

func IsFuncPkg(pass *analysis.Pass, fn *ast.CallExpr, pkgPath string) bool {
	pkg := GetFuncExprPkg(pass, fn.Fun)
	if pkg != nil {
		return pkg.Path() == pkgPath
	}
	return false
}

func GetFuncExprPkg(pass *analysis.Pass, e ast.Expr) *types.Package {
	ident := GetFuncExprIdent(pass, e)
	if ident == nil {
		return nil
	}
	return GetIdentPkg(pass, ident)
}

func GetIdentPkg(pass *analysis.Pass, ident *ast.Ident) *types.Package {
	obj := pass.TypesInfo.ObjectOf(ident)
	if obj == nil {
		return nil
	}
	return obj.Pkg()
}

func GetFuncExprObject(pass *analysis.Pass, fn *ast.CallExpr) types.Object {
	ident := GetFuncExprIdent(pass, fn.Fun)
	if ident == nil {
		return nil
	}
	obj := pass.TypesInfo.ObjectOf(ident)
	return obj
}

func IsBuiltinFunc(pass *analysis.Pass, fn *ast.CallExpr, name string) bool {
	obj := GetFuncExprObject(pass, fn)
	if obj == nil {
		return false
	}
	bo, ok := obj.(*types.Builtin)
	if !ok {
		return false
	}
	return bo.Name() == name
}

type FuncType struct {
	FuncName    string
	ArgsNum     int
	Signature   string
	ResultsNum  int
	Variadic    bool
	UseVariadic bool
}

func IsFuncType(pass *analysis.Pass, node *ast.CallExpr, fnType FuncType) bool {
	if len(node.Args) != fnType.ArgsNum {
		return false
	}

	if fnType.Variadic {
		useVariadic := node.Ellipsis != token.NoPos
		if fnType.UseVariadic != useVariadic {
			return false
		}
	}

	sign, ok := pass.TypesInfo.TypeOf(node.Fun).(*types.Signature)
	if !ok {
		return false
	}
	if sign.Variadic() != fnType.Variadic {
		return false
	}
	if fnType.Signature != "" {
		if sign.String() != fnType.Signature {
			return false
		}
	}
	if sign.Params().Len() != fnType.ArgsNum {
		return false
	}
	if sign.Results().Len() != fnType.ResultsNum {
		return false
	}
	return true
}

func GetCode(fset *token.FileSet, node any) (string, error) {
	buf := new(bytes.Buffer)
	if err := printer.Fprint(buf, fset, node); err != nil {
		return "", fmt.Errorf("unable to print node for %+v: %w", node, err)
	}
	return buf.String(), nil
}

func GetCodeSafe(fset *token.FileSet, node any) string {
	code, _ := GetCode(fset, node)
	return code
}

func IsBasicType(typ types.Type) bool {
	if basicType, ok := typ.(*types.Basic); ok {
		kind := basicType.Kind()
		switch kind {
		case types.Bool, types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64,
			types.Float32, types.Float64, types.Complex64, types.Complex128, types.String:
			return true
		}
	}
	return false
}

func IsTupleAll(tp *types.Tuple, predict func(types.Type) bool) bool {
	for i := 0; i < tp.Len(); i++ {
		if !predict(tp.At(i).Type()) {
			return false
		}
	}
	return true
}

func IsTestFile(pass *analysis.Pass, node ast.Expr) bool {
	return strings.HasSuffix(pass.Fset.Position(node.Pos()).Filename, "_test.go")
}

func IsBasicLit(node ast.Expr) bool {
	_, ok := node.(*ast.BasicLit)
	return ok
}

func IsCompositeLit(node ast.Expr) bool {
	return IsType[*ast.CompositeLit](node)
}

func IsType[T any](node ast.Expr) bool {
	_, ok := node.(T)
	return ok
}

func IsIdentNil(node ast.Expr) bool {
	v, ok := node.(*ast.Ident)
	return ok && v.Name == "nil" && v.Obj == nil
}

// restOfBlock, given a traversal stack, finds the innermost containing
// block and returns the suffix of its statements starting with the current
// node, along with the number of call expressions encountered.
func restOfBlock(stack []ast.Node) ([]ast.Stmt, int) {
	var ncalls int
	for i := len(stack) - 1; i >= 0; i-- {
		if b, ok := stack[i].(*ast.BlockStmt); ok {
			for j, v := range b.List {
				if v == stack[i+1] {
					return b.List[j:], ncalls
				}
			}
			break
		}

		if _, ok := stack[i].(*ast.CallExpr); ok {
			ncalls++
		}
	}
	return nil, 0
}
