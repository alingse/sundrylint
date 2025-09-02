package sundrylint

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func LintComparePtr(pass *analysis.Pass, node *ast.BinaryExpr) (ds []analysis.Diagnostic) {
	if node.Op != token.EQL && node.Op != token.NEQ {
		return nil
	}

	xType := pass.TypesInfo.TypeOf(node.X)
	yType := pass.TypesInfo.TypeOf(node.Y)

	if !isPointerType(xType) || !isPointerType(yType) {
		return nil
	}

	// Check if both pointers point to basic numeric types
	if !isBasicNumericPointerType(xType) || !isBasicNumericPointerType(yType) {
		return nil
	}

	// Check if both operands are struct fields with getters
	if !isStructFieldWithGetter(pass, node.X) || !isStructFieldWithGetter(pass, node.Y) {
		return nil
	}

	return []analysis.Diagnostic{
		{
			Pos:      node.Pos(),
			End:      node.End(),
			Category: LinterName,
			Message:  SubLinterComparePtrMessage,
		},
	}
}

func isPointerType(t types.Type) bool {
	_, ok := t.(*types.Pointer)
	return ok
}

func isBasicNumericPointerType(t types.Type) bool {
	ptrType, ok := t.(*types.Pointer)
	if !ok {
		return false
	}
	return isBasicNumericType(ptrType.Elem())
}

func isBasicNumericType(t types.Type) bool {
	basicType, ok := t.(*types.Basic)
	if !ok {
		return false
	}
	switch basicType.Kind() {
	case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
		types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64,
		types.Float32, types.Float64:
		return true
	default:
		return false
	}
}

func isStructFieldWithGetter(pass *analysis.Pass, expr ast.Expr) bool {
	// Check if the expression is a selector expression (struct.field)
	selExpr, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	// Get the type information for the selector
	selObj := pass.TypesInfo.ObjectOf(selExpr.Sel)
	if selObj == nil {
		return false
	}

	// Check if it's a struct field
	if _, ok := selObj.(*types.Var); !ok {
		return false
	}

	// Get the struct type
	structType := getStructType(pass, selExpr.X)
	if structType == nil {
		return false
	}

	// Check if the struct has a getter method for this field
	fieldName := selObj.Name()
	return hasGetterMethod(pass, selExpr.X, fieldName)
}

func getStructType(pass *analysis.Pass, expr ast.Expr) *types.Struct {
	// Get the type of the expression
	typ := pass.TypesInfo.TypeOf(expr)
	if typ == nil {
		return nil
	}

	// If it's a pointer, get the underlying type
	if ptrType, ok := typ.(*types.Pointer); ok {
		typ = ptrType.Elem()
	}

	// Check if it's a struct or named type
	switch t := typ.(type) {
	case *types.Struct:
		return t
	case *types.Named:
		if structType, ok := t.Underlying().(*types.Struct); ok {
			return structType
		}
		return nil
	default:
		return nil
	}
}

func hasGetterMethod(pass *analysis.Pass, expr ast.Expr, fieldName string) bool {
	// Construct the expected getter method name: Get + fieldName
	getterName := "Get" + fieldName

	// Get the type of the expression
	typ := pass.TypesInfo.TypeOf(expr)
	if typ == nil {
		return false
	}

	// If it's a pointer, get the underlying type
	if ptrType, ok := typ.(*types.Pointer); ok {
		typ = ptrType.Elem()
	}

	// Check if it's a named type
	namedType, ok := typ.(*types.Named)
	if !ok {
		return false
	}

	// Check methods directly on the named type
	for i := 0; i < namedType.NumMethods(); i++ {
		method := namedType.Method(i)
		if method.Name() == getterName {
			return isGetterMethodWithCorrectReturnType(method, "")
		}
	}

	return false
}

func isGetterMethodWithCorrectReturnType(method *types.Func, fieldName string) bool {
	sig := method.Type().(*types.Signature)

	// Check if the method has no parameters and one return value
	if sig.Params().Len() != 0 || sig.Results().Len() != 1 {
		return false
	}

	// Check if the return type is a basic numeric type
	returnType := sig.Results().At(0).Type()
	return isBasicNumericType(returnType)
}
