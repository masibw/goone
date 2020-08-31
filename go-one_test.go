package go-one_test

import (
	"testing"

	"github.com/masibw/go-one"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, go-one.Analyzer, "a")
}

