package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	s := flag.String("f", "sample.go", "file-path")
	flag.Parse()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, *s, nil, 0)
	if err != nil {
		fmt.Printf("Failedto parse file: cause %v", err)
		return
	}
	if err := ast.Print(fset, f); err != nil {
		fmt.Printf("Print failed: cause %v", err)
		return
	}
}
