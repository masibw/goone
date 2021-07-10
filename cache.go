package goone

import (
	"go/token"
	"golang.org/x/tools/go/analysis"
	"strconv"
	"sync"
)

type ReportCache struct {
	sync.Mutex
	reportMemo map[string]bool
}

func NewReportCache() *ReportCache {
	return &ReportCache{
		reportMemo: make(map[string]bool),
	}
}
func (m *ReportCache) toKey(pass *analysis.Pass, pos token.Pos) (key string) {
	posn := pass.Fset.Position(pos)
	fileName, lineNum := posn.Filename, posn.Line
	key = fileName + strconv.Itoa(lineNum)
	return
}
func (m *ReportCache) Set(pass *analysis.Pass, pos token.Pos, value bool) {
	key := m.toKey(pass, pos)
	m.reportMemo[key] = value
}

func (m *ReportCache) Get(pass *analysis.Pass, pos token.Pos) bool {
	key := m.toKey(pass, pos)
	value := m.reportMemo[key]
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


func (m *FuncCache) Get(key token.Pos) (bool,bool) {
	m.Lock()
	value, ok := m.funcMemo[key]
	m.Unlock()
	if !ok {
		return false, false
	}
	return value, true
}

type PkgCache struct {
	sync.Mutex
	pkgMemo map[string]bool
}

func NewPkgCache() *PkgCache {
	return &PkgCache{
		pkgMemo: make(map[string]bool),
	}
}

func (m *PkgCache) Set(name string, value bool) {
	m.pkgMemo[name] = value
}

func (m *PkgCache) Get(name string)  bool {
	value := m.pkgMemo[name]
	return value
}

func (m *PkgCache) Exists(key string) bool {
	m.Lock()
	_, exist := m.pkgMemo[key]
	m.Unlock()
	return exist
}

// pkgCache contains pkg AST
var pkgCache *PkgCache

// reportCache manages whether if this line already reported
var reportCache *ReportCache

// funcCache manages whether if this function contains queries
var funcCache *FuncCache