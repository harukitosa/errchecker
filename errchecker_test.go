package errchecker_test

import (
	"testing"

	"github.com/harukitosa/merucari/errchecker"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, errchecker.Analyzer, "a")
	analysistest.Run(t, testdata, errchecker.Analyzer, "b")
	analysistest.Run(t, testdata, errchecker.Analyzer, "sample_test")
}
