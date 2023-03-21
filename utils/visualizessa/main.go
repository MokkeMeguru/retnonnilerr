package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"

	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

func main() {
	s := flag.String("f", "sample.go", "file-path")
	d := flag.String("d", "./", "directory-path")
	flag.Parse()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, *s, nil, parser.AllErrors)
	if err != nil {
		log.Fatalf("Failed to parse file at parser: cause %v", err)
	}
	pkgs, err := parser.ParseDir(fset, *d, nil, parser.AllErrors)
	if err != nil {
		log.Fatalf("Failed to parse dir at parser: cause %v", err)
	}
	files := []*ast.File{}
	for _, pkg := range pkgs {
		if pkg.Name != f.Name.Name {
			continue
		}
		for _, _f := range pkg.Files {
			files = append(files, _f)
		}
	}
	typesPkg := types.NewPackage(f.Name.Name, "")
	ssaPkg, _, err := ssautil.BuildPackage(&types.Config{
		Importer: importer.Default(),
	}, fset, typesPkg, files, ssa.GlobalDebug)
	if err != nil {
		log.Fatalf("Failed to parse file at ssa: cause %v", err)
	}
	ssaPkg.WriteTo(os.Stdout)
	for _, member := range ssaPkg.Members {
		if fn, ok := member.(*ssa.Function); ok {
			if fmt.Sprintf("./%s", fset.Position(fn.Pos()).Filename) != fset.Position(f.Pos()).Filename {
				continue
			}
			ssaPkg.Func(fn.Name()).WriteTo(os.Stdout)
		}
	}
}
