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
	goone.Analyzer.Flags.Set("configPath", defaultPath)
	analysistest.Run(t, testdata, goone.Analyzer, "base", "separated", "gorm", "gorp", "sqlx", "user_def")
}
