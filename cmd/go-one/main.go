package main

import (
	"github.com/masibw/go-one"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(go-one.Analyzer) }

