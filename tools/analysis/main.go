package main

import (
	"github.com/MeguruMokke/retnonnilerr"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		retnonnilerr.Analyzer,
	)
}
