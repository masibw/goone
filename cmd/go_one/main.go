package main

import (
	"github.com/masibw/goone"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(goone.Analyzer) }
