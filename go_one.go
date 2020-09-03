package go_one

import (
	_ "database/sql"
	"fmt"
	"github.com/gostaticanalysis/analysisutil"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/ast/inspector"
	"log"
)

const doc = "go_one finds N+1 query "

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "go_one",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	forFilter := []ast.Node{
		(*ast.ForStmt)(nil),
		(*ast.RangeStmt)(nil),
	}

	inspect.Preorder(forFilter, func(n ast.Node) {
		//s := pass.Pkg.Scope()
		//for _, name:=range s.Names(){
		//	obj := s.Lookup(name)
		//	objs :=pass.TypesInfo.Defs[obj]
		//	fmt.Println(obj.Name(),	obj.Type().Underlying())
		//}

		switch n := n.(type) {
		case *ast.ForStmt, *ast.RangeStmt:
			ast.Inspect(n, func(n ast.Node) bool {
				switch node := n.(type) {
				case *ast.Ident:
					if tv, ok := pass.TypesInfo.Types[node]; ok {
						obj := analysisutil.TypeOf(pass, "database/sql", "*DB")
						if types.Identical(tv.Type, obj) {
							pass.Reportf(node.Pos(), "this query might be causes bad performance")
							log.Printf("%s is detected %d ", node.Name, node.NamePos)
						}
					}
				case *ast.CallExpr:
					switch funcExpr := node.Fun.(type) {
					case *ast.Ident:
						//fmt.Println(funcExpr.Name)
						obj := funcExpr.Obj
						if obj == nil {
							//if anotherFileNode := analysisutil.TypeOf(pass,funcExpr.Name); anotherFileNode != nil {
							//
							//}
							if anotherFileNode := pass.TypesInfo.ObjectOf(funcExpr); anotherFileNode != nil {
								file := analysisutil.File(pass, anotherFileNode.Pos())
								if file == nil {
									return false
								}
								if path, ok := astutil.PathEnclosingInterval(file, anotherFileNode.Pos(), anotherFileNode.Pos()); ok {
									if funcDecl, ok := path[1].(*ast.FuncDecl); ok {
										fmt.Println("funcName", funcDecl.Name)
										ast.Inspect(funcDecl, func(n ast.Node) bool {
											switch ident := n.(type) {
											case *ast.Ident:
												if tv, ok := pass.TypesInfo.Types[ident]; ok {
													obj := analysisutil.TypeOf(pass, "database/sql", "*DB")
													if types.Identical(tv.Type, obj) {
														pass.Reportf(node.Pos(), "this query might be causes bad performance")
														log.Printf("%s is detected %d ", ident.Name, ident.NamePos)
													}
												}

											}
											return true
										})
									}
								}

								//fmt.Println("anotherFileNode",anotherFileNode.Name())
							}

							return false
						}
						switch decl := obj.Decl.(type) {
						case *ast.FuncDecl:
							ast.Inspect(decl, func(n ast.Node) bool {
								switch ident := n.(type) {
								case *ast.Ident:
									if tv, ok := pass.TypesInfo.Types[ident]; ok {
										obj := analysisutil.TypeOf(pass, "database/sql", "*DB")
										if types.Identical(tv.Type, obj) {
											pass.Reportf(node.Pos(), "this query might be causes bad performance")
											log.Printf("%s is detected %d ", ident.Name, ident.NamePos)
										}
									}

								}
								return true
							})
						}

					}

				}
				return true

			})

		}

	})

	return nil, nil
}
