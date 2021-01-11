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
	defaultPath, err := filepath.Abs("goone.yml")
	if err != nil {
		log.Println(err)
	}
	err = goone.Analyzer.Flags.Set("configPath", defaultPath)
	if err != nil {
		log.Println(err)
	}
	analysistest.Run(t, testdata, goone.Analyzer, "another_package", "another_package_dot", "base", "dummy_type", "gorm", "gorp", "separated", "sqlx", "user_def")
}
