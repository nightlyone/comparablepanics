// Package comparablepanics reports generic calls that may panic when using
// the comparable constraint under Go 1.20+.
package comparablepanics

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer reports generic calls whose comparable type arguments may panic at runtime
// on Go 1.20+.
var Analyzer = &analysis.Analyzer{
	Name: "comparablepanics",
	Doc:  Doc,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		usesgenerics.Analyzer,
	},
	Run: run,
}

// Doc is the analyzer description shown by go doc and tooling.
const Doc = `detects whether a comparable type instantiation can panic at runtime
so code written for Go 1.18 or 1.19 under the assumption that comparable cannot panic
will also not panic in Go 1.20 and later`

// run checks the AST for generic function instantiations and reports any
// comparable-constraint arguments that may panic at runtime.
func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	hasGenerics := pass.ResultOf[usesgenerics.Analyzer].(*usesgenerics.Result)

	// No generics? Ignore.
	if hasGenerics.Direct&(usesgenerics.FuncInstantiation|usesgenerics.TypeInstantiation) == 0 {
		return nil, nil
	}
	funcsWithGenerics := make(map[string]types.Type)
	inspect.Preorder(nil, func(n ast.Node) {
		runFunc(pass, n, funcsWithGenerics)
	})
	return nil, nil
}

// runFunc records generic function declarations and dispatches calls for
// comparable-constraint checks.
func runFunc(pass *analysis.Pass, node ast.Node, funcsWithGenerics map[string]types.Type) {
	switch v := node.(type) {
	case *ast.CallExpr:
		runCall(pass, v, funcsWithGenerics)
	case *ast.FuncDecl:
		if v.Type.TypeParams == nil {
			return
		}
		funcsWithGenerics[v.Name.Name] = pass.TypesInfo.TypeOf(v.Name)
	}

}

// runCall inspects a call expression and reports arguments that are not safe for
// a comparable type parameter under Go 1.20+.
func runCall(pass *analysis.Pass, call *ast.CallExpr, funcsWithGenerics map[string]types.Type) {
	if len(call.Args) == 0 {
		return
	}
	fun := call.Fun
	var funName string
	if id, ok := fun.(*ast.Ident); ok {
		funName = id.Name
	}
	typ := funcsWithGenerics[funName]
	if typ == nil {
		return
	}
	sig, ok := typ.(*types.Signature)
	if !ok {
		return
	}
	if sig.TypeParams().Len() == 0 {
		return
	}
	comparableTypeNames := make(map[string]struct{})
	for tp := range sig.TypeParams().TypeParams() {
		if !isComparableTypeName(tp) {
			continue
		}
		comparableTypeNames[tp.Obj().Id()] = struct{}{}
	}
	if len(comparableTypeNames) == 0 {
		return
	}
	for i, arg := range call.Args {
		if !isComparableWithoutPanic(arg, pass.TypesInfo) {
			if ident, ok := arg.(*ast.Ident); ok {
				pass.Reportf(arg.Pos(), "argument %d/%d (%s) constrained as comparable may panic at runtime in Go 1.20+", i+1, len(call.Args), ident.Name)
				continue
			}
			pass.Reportf(arg.Pos(), "argument %d/%d constrained as comparable may panic at runtime in Go 1.20+", i+1, len(call.Args))
		}
	}

}

// isComparableWithoutPanic reports whether an expression's type is comparable
// without triggering the Go 1.20+ panic case for a comparable constraint.
func isComparableWithoutPanic(arg ast.Expr, info *types.Info) bool {
	typ := info.TypeOf(arg)
	var recursiveIsComparableWithoutPanic func(typ types.Type) bool
	recursiveIsComparableWithoutPanic = func(typ types.Type) bool {
		if typ == nil {
			return false
		}
		if !types.Comparable(typ) {
			return false
		}

		typ = typ.Underlying()
		if _, ok := typ.(*types.Basic); ok {
			return true
		}
		if st, ok := typ.(*types.Struct); ok {
			for field := range st.Fields() {
				if !recursiveIsComparableWithoutPanic(field.Type()) {
					return false
				}
			}
			return true
		}
		if at, ok := typ.(*types.Array); ok {
			return recursiveIsComparableWithoutPanic(at.Elem())
		}
		return false
	}
	return recursiveIsComparableWithoutPanic(typ)
}

// isComparableTypeName reports whether the type parameter is constrained by
// the comparable interface or by a type parameter named comparable.
func isComparableTypeName(tp *types.TypeParam) bool {
	if tp.Obj().Name() == "comparable" && tp.Obj().Pkg() == nil {
		return true
	}
	return tp.Constraint().String() == "comparable"
}
