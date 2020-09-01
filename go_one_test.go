package go_one_test

import (
	go_one "github.com/masibw/go_one"
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, go_one.Analyzer,  "separated")
}
