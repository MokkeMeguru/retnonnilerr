package main

import (
	"github.com/MeguruMokke/retnonnilerr"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(retnonnilerr.Analyzer)
}
