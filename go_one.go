package go_one

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "go-one is ..."

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

	nodeFilter := []ast.Node{
		(*ast.ForStmt)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
				switch n:=n.(type){
				case *ast.ForStmt:

					for _, stmt := range n.Body.List{
						switch stmt := stmt.(type){
						case *ast.AssignStmt:
						for _,expr := range stmt.Rhs{
							switch expr := expr.(type){
							case *ast.CallExpr:
								switch fun:=expr.Fun.(type){
								case *ast.SelectorExpr:
									//ast.Print(nil,fun)
									switch x:=fun.X.(type) {
									case *ast.CallExpr:
										switch fun := x.Fun.(type) {
										case *ast.SelectorExpr:
											switch x := fun.X.(type) {
											case *ast.Ident:
												if tv, ok := pass.TypesInfo.Types[x]; ok {
													if tv.Type.String() == "*database/sql.DB" {
														pass.Reportf(fun.Pos(), "this query might be causes bad performance")
													}
												}
											}
										}
									}
								}
							}
						}
						}
					}
				}

	})

	return nil, nil
}

