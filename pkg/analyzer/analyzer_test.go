package analyzer

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"os"
	"path/filepath"
	"testing"
)

func TestAll(t *testing.T) {
	wd, _ := os.Getwd()
	testdata := filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata")
	Analyzer.Flags.Set("path", filepath.Join(filepath.Dir(wd), "localize", "locales"))
	analysistest.Run(t, testdata, Analyzer, "p")
}
