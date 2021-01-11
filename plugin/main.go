// This file can build as base plugin for golangci-lint by below command.
//    go build -buildmode=plugin -o path_to_plugin_dir github.com/masibw/goone/plugin/goone
// See: https://golangci-lint.run/contributing/new-linters/#how-to-add-a-private-linter-to-golangci-lint

package main

import (
	"strings"

	"github.com/masibw/goone"
	"golang.org/x/tools/go/analysis"
)

// flags for Analyzer.Flag.
// If you would like to specify flags for your plugin, you can put them via 'ldflags' as below.
//     $ go build -buildmode=plugin -ldflags "-X 'main.flags=-opt val'" github.com/masibw/goone/plugin/goone
var flags string

// AnalyzerPlugin provides analyzers as base plugin.
// It follows golangci-lint style plugin.
var AnalyzerPlugin analyzerPlugin //nolint

type analyzerPlugin struct{}

func (analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	if flags != "" {
		flagset := goone.Analyzer.Flags
		if err := flagset.Parse(strings.Split(flags, " ")); err != nil {
			panic("cannot parse flags of goone: " + err.Error())
		}
	}
	return []*analysis.Analyzer{
		goone.Analyzer,
	}
}
