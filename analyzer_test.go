package comparablepanics_test

import (
	"testing"

	"github.com/nightlyone/comparablepanics"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, comparablepanics.Analyzer,
		"a",
		"b",
		"c",
	)
}
