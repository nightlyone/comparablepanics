// Package usesgenerics defines an Analyzer that checks for usage of generic
// features added in Go 1.18.
package comparablepanics

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "comparablepanics",
	Doc:  Doc,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		usesgenerics.Analyzer,
	},
	Run: run,
}

const Doc = `detects whether a comparable type instantiations can panic at runtime
to ensure code written for Go 1.18 or 1.19 under the assumption that comparable cannot panic
will also not panic in Go 1.20 and later`

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	hasGenerics := pass.ResultOf[usesgenerics.Analyzer].(*usesgenerics.Result)

	// No generics? Ignore.
	if hasGenerics.Direct&(usesgenerics.FuncInstantiation|usesgenerics.TypeInstantiation) == 0 {
		fmt.Println("no generics")
		return nil, nil
	}
	funcsWithGenerics := make(map[string]types.Type)
	inspect.Preorder(nil, func(n ast.Node) {
		runFunc(pass, n, funcsWithGenerics)
	})
	return nil, nil
}

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
	for i := 0; i < sig.TypeParams().Len(); i++ {
		tp := sig.TypeParams().At(i)
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
			for i := 0; i < st.NumFields(); i++ {
				if !recursiveIsComparableWithoutPanic(st.Field(i).Type()) {
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

func isComparableTypeName(tp *types.TypeParam) bool {
	if tp.Obj().Name() == "comparable" && tp.Obj().Pkg() == nil {
		return true
	}
	return tp.Constraint().String() == "comparable"
}
