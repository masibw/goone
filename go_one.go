package goone

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io/ioutil"
	"log"
	"strings"
	"sync"

	"golang.org/x/tools/go/packages"

	"gopkg.in/yaml.v2"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

type Types struct {
	Package []struct {
		PkgName   string `yaml:"pkgName"`
		TypeNames []struct {
			TypeName string `yaml:"typeName"`
		} `yaml:"typeNames"`
	} `yaml:"package"`
}

const doc = "go_one finds N+1 query "

type SearchCache struct {
	sync.Mutex
	searchMemo map[token.Pos]bool
}

func NewSearchCache() *SearchCache {
	return &SearchCache{
		searchMemo: make(map[token.Pos]bool),
	}
}

func (m *SearchCache) Set(key token.Pos, value bool) {
	m.Lock()
	m.searchMemo[key] = value
	m.Unlock()
}

func (m *SearchCache) Get(key token.Pos) bool {
	m.Lock()
	value := m.searchMemo[key]
	m.Unlock()
	return value
}

type FuncCache struct {
	sync.Mutex
	funcMemo map[token.Pos]bool
}

func NewFuncCache() *FuncCache {
	return &FuncCache{
		funcMemo: make(map[token.Pos]bool),
	}
}

func (m *FuncCache) Set(key token.Pos, value bool) {
	m.Lock()
	m.funcMemo[key] = value
	m.Unlock()
}

func (m *FuncCache) Exists(key token.Pos) bool {
	m.Lock()
	_, exist := m.funcMemo[key]
	m.Unlock()
	return exist
}

func (m *FuncCache) Get(key token.Pos) bool {
	m.Lock()
	value := m.funcMemo[key]
	m.Unlock()
	return value
}

// searchCache manages whether if already searched this node
var searchCache *SearchCache

// funcCache manages whether if this function contains queries
var funcCache *FuncCache
var configPath string

// Analyzer is analysis files
var Analyzer = &analysis.Analyzer{
	Name: "go_one",
	Doc:  doc,
	Run:  run,
	//FactTypes: []analysis.Fact{new(isWrapper)}, // When Fact is specified, Analyzer also runs on imported files(packages.Loadで読み取るのでいらないはず)
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}
var sqlTypes []types.Type

type isWrapper struct{}

func (f *isWrapper) AFact() {}

func init() {
	Analyzer.Flags.StringVar(&configPath, "configPath", "", "config file path(abs)")
}

func appendTypes(pass *analysis.Pass, pkg, name string) {
	if typ := analysisutil.TypeOf(pass, pkg, name); typ != nil {
		sqlTypes = append(sqlTypes, typ)
	}
}

func readTypeConfig(configPath string) *Types {
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil
	}
	typesFromConfig := Types{}
	err = yaml.Unmarshal([]byte(buf), &typesFromConfig)
	if err != nil {
		log.Fatalf("yml parse error:%v", err)
	}

	return &typesFromConfig
}

func prepareTypes(pass *analysis.Pass, configPath string) {
	typesFromConfig := readTypeConfig(configPath)
	if typesFromConfig != nil {
		for _, pkg := range typesFromConfig.Package {
			pkgName := pkg.PkgName
			for _, typeName := range pkg.TypeNames {
				appendTypes(pass, pkgName, typeName.TypeName)
			}
		}
	}

	appendTypes(pass, "database/sql", "*DB")
	appendTypes(pass, "gorm.io/gorm", "*DB")
	appendTypes(pass, "gopkg.in/gorp.v1", "*DbMap")
	appendTypes(pass, "gopkg.in/gorp.v2", "*DbMap")
	appendTypes(pass, "github.com/go-gorp/gorp/v3", "*DbMap")
	appendTypes(pass, "github.com/jmoiron/sqlx", "*DB")
}

func run(pass *analysis.Pass) (interface{}, error) {
	searchCache = NewSearchCache()
	funcCache = NewFuncCache()
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	prepareTypes(pass, configPath)

	if sqlTypes == nil {
		return nil, nil
	}
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
	return
}

func anotherFileSearch(pass *analysis.Pass, funcExpr *ast.Ident, parentNode ast.Node) bool {
	if anotherFileNode := pass.TypesInfo.ObjectOf(funcExpr); anotherFileNode != nil {
		file := analysisutil.File(pass, anotherFileNode.Pos())
		if file == nil {
			if anotherFileNode.Pkg() != nil {
				importPath := convertToImportPath(pass, fmt.Sprintf("%s", anotherFileNode.Pkg().Name()))
				pkgs := loadImportPackages(importPath)

				for _, pkg := range pkgs {
					for _, file := range pkg.Syntax {
						inspectFile(pass, parentNode, file, funcExpr.Name, pkg.TypesInfo)
					}
				}
			}
			return false
		}
		inspectFile(pass, parentNode, file, funcExpr.Name, nil)
	}

	return false
}

func findQuery(pass *analysis.Pass, rootNode, parentNode ast.Node, pkgTypes *types.Info) {

	if funcCache.Exists(rootNode.Pos()) {
		if funcCache.Get(rootNode.Pos()) {
			pass.Reportf(parentNode.Pos(), "this query is called in a loop")
		}
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

				for _, typ := range sqlTypes {
					//TODO Comparing by string is bad, but I don't have any ideas to compare another package's type
					if tv.Type.String() == typ.String() {
						//if types.Identical(tv.Type,typ){
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
					if !searchCache.Get(decl.Pos()) {
						searchCache.Set(decl.Pos(), true)
						newParentNode := parentNode
						if parentNode == nil {
							newParentNode = node
						}

						findQuery(pass, decl, newParentNode, nil)
					} else {
						if funcCache.Get(decl.Pos()) {
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

					pkgs := loadImportPackages(importPath)
					for _, pkg := range pkgs {

						//Do not scan standard packages. Is it...ok?
						if !strings.Contains(pkg.PkgPath, ".") {
							continue
						}
						for _, file := range pkg.Syntax {
							newParentNode := parentNode
							if parentNode == nil {
								newParentNode = node
							}

							inspectFile(pass, newParentNode, file, funcExpr.Sel.Name, pkg.TypesInfo)
						}
					}
				}

			}

		}
		return true

	})
}

func convertToImportPath(pass *analysis.Pass, pkgName string) (importPath string) {
	for _, v := range pass.Pkg.Imports() {
		if strings.HasSuffix("/"+v.Path(), pkgName) {
			importPath = v.Path()
			if strings.HasPrefix(importPath, "vendor/") {
				importPath = strings.TrimPrefix(importPath, "vendor/")
			}
			return
		}
	}
	return
}
