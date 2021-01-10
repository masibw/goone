package goone_test

import (
	"log"
	"path/filepath"
	"testing"

	"github.com/masibw/goone"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is base test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	defaultPath, err := filepath.Abs("go_one.yml")
	if err != nil {
		log.Println(err)
	}
	err = goone.Analyzer.Flags.Set("configPath", defaultPath)
	if err != nil {
		log.Println(err)
	}
	analysistest.Run(t, testdata, goone.Analyzer, "separated")
}
