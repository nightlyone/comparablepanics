// The comparablepanics command applies the
// [github.com/nightlyone/comparablepanics]
// analysis to the specified packages of Go source code.
package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/nightlyone/comparablepanics"
)

func main() { singlechecker.Main(comparablepanics.Analyzer) }
