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

func helper(pass *analysis.Pass, node ast.Node){

}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	forFilter := []ast.Node{
		(*ast.ForStmt)(nil),
	}

	inspect.Preorder(forFilter, func(n ast.Node) {
				switch n:=n.(type){
				case *ast.ForStmt:
					ast.Inspect(n,func(n ast.Node) bool{
						switch node := n.(type){
						case *ast.Ident:
							if tv, ok := pass.TypesInfo.Types[node]; ok {
								if tv.Type.String() == "*database/sql.DB" {
									pass.Reportf(node.Pos(), "this query might be causes bad performance")
									return false
								}
							}
						}
						return true
					})
				}

	})

	return nil, nil
}

