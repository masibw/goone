package goone

import (
	"fmt"
	"github.com/gostaticanalysis/analysisutil"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
	"log"
	"strings"
)

type Types struct {
	Package []struct {
		PkgName   string `yaml:"pkgName"`
		TypeNames []struct {
			TypeName string `yaml:"typeName"`
		} `yaml:"typeNames"`
	} `yaml:"package"`
}

const doc = "goone finds N+1 query "

var configPath string

// Analyzer is analysis files
var Analyzer = &analysis.Analyzer{
	Name: "goone",
	Doc:  doc,
	Run:  run,
	//FactTypes: []analysis.Fact{new(isWrapper)}, // When Fact is specified, Analyzer also runs on imported files(packages.Loadで読み取るのでいらないはず)
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}
var sqlTypes []string

func run(pass *analysis.Pass) (interface{}, error) {
	funcCache = NewFuncCache()
	reportCache = NewReportCache()
	pkgCache = NewPkgCache()
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	prepareTypes(pass)
	forFilter := []ast.Node{
		(*ast.ForStmt)(nil),
		(*ast.RangeStmt)(nil),
	}

	inspect.Preorder(forFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.ForStmt, *ast.RangeStmt:
				findQuery(pass, n, nil, nil)
		}
	})
	return nil, nil
}

func inspectFile(pass *analysis.Pass, parentNode ast.Node, file *ast.File, funcName string, pkgTypes *types.Info) {
	inspect := inspector.New([]*ast.File{file})
	types := []ast.Node{new(ast.FuncDecl)}
	inspect.Preorder(types, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			if n.Name.Name == funcName {
					findQuery(pass, n, parentNode, pkgTypes)
			}
		}
	})
}

func loadImportPackages(pkgName string) (pkgs []*packages.Package) {
	config := &packages.Config{Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedSyntax | packages.NeedTypesInfo}
	var err error
	pkgs, err = packages.Load(config, pkgName)
	if err != nil {
		log.Fatalf(err.Error())
	}
	pkgCache.Set(pkgName,true)
	return
}

func anotherFileSearch(pass *analysis.Pass, funcExpr *ast.Ident, parentNode ast.Node) bool {
	if anotherFileNode := pass.TypesInfo.ObjectOf(funcExpr); anotherFileNode != nil {
		file := analysisutil.File(pass, anotherFileNode.Pos())
		if file == nil {
			if anotherFileNode.Pkg() != nil {
				importPath := convertToImportPath(pass, anotherFileNode.Pkg().Name())

				if pkgCache.Exists(importPath){
					// Check whether if already checked function
					if containsQuery, ok := funcCache.Get(funcExpr.Pos()); ok{
						if containsQuery {
							pass.Reportf(parentNode.Pos(), "this query is called in a loop")
						}
						return false
					}
				}
				// If the package has not been loaded yet, load it.
				pkgs := loadImportPackages(importPath)
				for i := range pkgs {
					for j := range pkgs[i].Syntax {
							inspectFile(pass, parentNode, pkgs[i].Syntax[j], funcExpr.Name, pkgs[i].TypesInfo)
					}
				}
			}
			return false
		}
		inspectFile(pass, parentNode, file, funcExpr.Name, nil)
	}

	return false
}

func findQuery(pass *analysis.Pass, rootNode, parentNode ast.Node, pkgTypes *types.Info) { //nolint:gocognit

	if containsQuery, ok := funcCache.Get(rootNode.Pos()); ok && containsQuery{

			reportCache.Lock()
			if reportCache.Get(pass, parentNode.Pos()) {
				reportCache.Unlock()
				return
			}
			reportCache.Set(pass, parentNode.Pos(), true)
			reportCache.Unlock()

			pass.Reportf(parentNode.Pos(), "this query is called in a loop")

		return
	}

	ast.Inspect(rootNode, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.Ident:

			// pass doesn't have a separate package typesInfo, so we need to pass the pkg info.
			var tv types.TypeAndValue
			var ok bool

			// Use types of pass(same package) or pkgTypes(another package)
			if tv, ok = pass.TypesInfo.Types[node]; !ok && pkgTypes != nil {
				if tvTmp, exist := pkgTypes.Types[node]; exist {
					tv, ok = tvTmp, exist
				}
			}
			if ok {
				reportNode := parentNode
				if reportNode == nil {
					reportNode = node
				}
				for i := range sqlTypes {
					//TODO Comparing by string is bad, but I don't have any ideas to compare another package's type
					if strings.TrimPrefix(tv.Type.String(), "vendor/") == sqlTypes[i] || strings.TrimPrefix(tv.Type.String(), "*vendor/") == strings.TrimPrefix(sqlTypes[i], "*") {
						//if types.Identical(tv.Type,typ){
						reportCache.Lock()
						if reportCache.Get(pass, reportNode.Pos()) {
							reportCache.Unlock()
							return false
						}
						reportCache.Set(pass, reportNode.Pos(), true)
						reportCache.Unlock()

						pass.Reportf(reportNode.Pos(), "this query is called in a loop")
						funcCache.Set(rootNode.Pos(), true)

						return false
					}
				}

			}
		case *ast.CallExpr:
			switch funcExpr := node.Fun.(type) {
			case *ast.Ident:
				obj := funcExpr.Obj
				//if function does not exist in same file
				if obj == nil {
					return anotherFileSearch(pass, funcExpr, node)
				}
				switch decl := obj.Decl.(type) {
				case *ast.FuncDecl:
					if _, ok := funcCache.Get(decl.Pos()); !ok {
						newParentNode := parentNode
						if parentNode == nil {
							newParentNode = node
						}
							findQuery(pass, decl, newParentNode, nil)
					} else {
						if containsQuery,_ := funcCache.Get(decl.Pos()); containsQuery{

							reportCache.Lock()
							if reportCache.Get(pass, node.Pos()) {
								reportCache.Unlock()
								return false
							}
							reportCache.Set(pass, node.Pos(), true)
							reportCache.Unlock()

							pass.Reportf(node.Pos(), "this query is called in a loop")
						}
					}

				}
			//inspect another package file
			case *ast.SelectorExpr:
				obj := funcExpr.Sel.Obj
				if obj == nil {
					//TODO get packageName without fmt.Sprinf
					importPath := convertToImportPath(pass, fmt.Sprintf("%s", funcExpr.X))

					if pkgCache.Exists(importPath){
						if containsQuery, ok := funcCache.Get(funcExpr.Pos()); ok{
							if containsQuery {
								pass.Reportf(parentNode.Pos(), "this query is called in a loop")
							}
							return false
						}
					}
					// If the package has not been loaded yet, load it
					pkgs := loadImportPackages(importPath)
					for i := range pkgs {

						//Do not scan standard packages. Is it...ok?
						if !strings.Contains(pkgs[i].PkgPath, ".") {
							continue
						}
						for j := range pkgs[i].Syntax {
							newParentNode := parentNode
							if parentNode == nil {
								newParentNode = node
							}
							inspectFile(pass, newParentNode, pkgs[i].Syntax[j], funcExpr.Sel.Name, pkgs[i].TypesInfo)
						}
					}
				}

			}

		}
		return true

	})
}

func convertToImportPath(pass *analysis.Pass, pkgName string) (importPath string) {
	for i := range pass.Pkg.Imports() {
		if strings.HasSuffix("/"+pass.Pkg.Imports()[i].Path(), pkgName) {
			importPath = pass.Pkg.Imports()[i].Path()
			if strings.HasPrefix(importPath, "vendor/") {
				importPath = strings.TrimPrefix(importPath, "vendor/")
			}
			return
		}
	}
	return
}
