package main

import (
	"github.com/MokkeMeguru/retnonnilerr"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		retnonnilerr.Analyzer,
	)
}
