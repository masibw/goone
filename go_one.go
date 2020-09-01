package go_one

import (
	_ "database/sql"
	"github.com/gostaticanalysis/analysisutil"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
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
	}

	inspect.Preorder(forFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.ForStmt:
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
						//ast.Print(nil,funcExpr.Obj)
						obj := funcExpr.Obj
						if obj == nil {
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
