package go_one

import (
	_ "database/sql"
	"github.com/gostaticanalysis/analysisutil"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/ast/inspector"
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
		switch n := n.(type) {
		case *ast.ForStmt, *ast.RangeStmt:
			findQuery(pass, n, nil)
		}

	})

	return nil, nil
}

func findQuery(pass *analysis.Pass, rootNode, parentNode ast.Node) {
	ast.Inspect(rootNode, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.Ident:
			if tv, ok := pass.TypesInfo.Types[node]; ok {
				obj := analysisutil.TypeOf(pass, "database/sql", "*DB")
				if types.Identical(tv.Type, obj) {
					reportNode := parentNode
					if reportNode == nil {
						reportNode = node
					}
					pass.Reportf(reportNode.Pos(), "this query might be causes bad performance")
				}
			}
		case *ast.CallExpr:
			switch funcExpr := node.Fun.(type) {
			case *ast.Ident:
				obj := funcExpr.Obj
				//if function does not exist in same file
				if obj == nil {
					if anotherFileNode := pass.TypesInfo.ObjectOf(funcExpr); anotherFileNode != nil {
						file := analysisutil.File(pass, anotherFileNode.Pos())
						if file == nil {
							return false
						}
						if path, ok := astutil.PathEnclosingInterval(file, anotherFileNode.Pos(), anotherFileNode.Pos()); ok {
							if funcDecl, ok := path[1].(*ast.FuncDecl); ok {
								findQuery(pass, funcDecl, node)
							}
						}
					}

					return false
				}
				//if function exists in same file
				switch decl := obj.Decl.(type) {
				case *ast.FuncDecl:
					findQuery(pass, decl, node)
				}

			}

		}
		return true

	})
}
