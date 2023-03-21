package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	s := flag.String("f", "sample.go", "file-path")
	flag.Parse()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, *s, nil, 0)
	if err != nil {
		log.Fatalf("Failed to parse file: cause %v", err)
	}
	if err := ast.Print(fset, f); err != nil {
		log.Fatalf("Print failed: cause %v", err)
	}
}
