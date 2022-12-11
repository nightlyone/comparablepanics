package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/nightlyone/comparablepanics"
)

func main() { singlechecker.Main(comparablepanics.Analyzer) }
