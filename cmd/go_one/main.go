package main

import (
	"github.com/masibw/go_one"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(go_one.Analyzer) }
