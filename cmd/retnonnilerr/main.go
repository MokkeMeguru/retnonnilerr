package main

import (
	"github.com/MokkeMeguru/retnonnilerr"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(retnonnilerr.Analyzer)
}
