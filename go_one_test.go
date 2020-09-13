package goone_test

import (
	"testing"

	"github.com/masibw/goone"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is base test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, goone.Analyzer, "base", "separated", "gorm", "gorp", "sqlx")
}
