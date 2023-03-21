package retnonnilerr_test

import (
	"testing"

	"github.com/MokkeMeguru/retnonnilerr"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test_Run(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, retnonnilerr.Analyzer, "a")
}
