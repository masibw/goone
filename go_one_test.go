package go_one_test

import (
	"github.com/masibw/go_one"
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

// TestAnalyzer is base test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, go_one.Analyzer, "base", "separated")
}
