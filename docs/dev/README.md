## Project initialize

do

- `go install github.com/gostaticanalysis/skeleton/v2@latest`
- `skeleton -kind=ssa github.com/retnonnilerr`

references:

- [【Go 言語】自作した静的解析ツールを GitHub Actions で実行する](https://zenn.dev/ytakaya/articles/55a07808c2fd5e)

```
.
├── Makefile
├── README.md
├── cmd // for cli
│   └── {package-name}
│       └── main.go
├── go.mod
├── {package-name}.go // body
├── {package-name}_test.go
├── staticcheck.conf // ref: https://staticcheck.io/docs/configuration/
└── testdata
    └── src
        └── a
            ├── a.go
            └── a_test.go

```

### Initial main code

The first code for static analyzer is following one.

When we need to check the return type, we can use the SSA (Static Single Assignment) form.

```golang
package retnonnilerr

import (
	"fmt"
	"go/types"

	"github.com/gostaticanalysis/comment/passes/commentmap"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const (
	name = "retnonnilerr"
	doc  = "Retnonnilerr is a static analysis tool to detect `if err != nil { return nil }`"
)

var Analyzer = &analysis.Analyzer{
	Name: name,
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
		inspect.Analyzer,
		commentmap.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	s := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	for _, f := range s.SrcFuncs {
		fmt.Println(f)
		for _, b := range f.Blocks {
			fmt.Printf("\tBlock %d\n", b.Index)
			for _, instr := range b.Instrs {
				fmt.Printf("\t\t%[1]T\t%[1]v\n", instr)
				for _, v := range instr.Operands(nil) {
					if v != nil {
						fmt.Printf("\t\t\t%[1]T\t%[1]v\n", *v)
					}
				}
			}
		}
	}
	return nil, nil
}
```

and also, the test code for the analyzer is this.

```golang
package retnonnilerr_test

import (
	"testing"

	"github.com/MokkeMeguru/retnonnilerr"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test_Run(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, retnonnilerr.Analyzer, "a")
}
```

I recommend you add the Makefile as following to test our analyzer easier.

```makefile
deps:
	go install github.com/kisielk/errcheck@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

lint: deps
	go vet ./...
	errcheck ./...
	staticcheck ./...

test:
	go test $(go list ./... | grep -v /testdata)
```

## Add the first test data

Let's write the first test code. (our analyzer should pass these tests)

```golang
package a

type T struct {
	I int
}

func funcA() (*T, error) {
	return nil, nil
}
```

Test it.

```
❯ make test
go test
PASS
ok      github.com/MokkeMeguru/retnonnilerr     1.349s
```

And then, add the red (failed) test code.

```golang
func funcB() (*T, error) {
	var err error
	if err != nil {
		return nil, nil // want "return err"
	}
	return nil, nil
}
```

Test it.

```
❯ make test
go test
--- FAIL: Test_Run (0.92s)
    analysistest.go:520: a/a.go:14: no diagnostic was reported matching `return err`
FAIL
exit status 1
FAIL    github.com/MokkeMeguru/retnonnilerr     1.087s
make: *** [test] Error 1
```

We failed the tests. That means, we are in the TDD's (red-green-refactor) process now.

## write the test's green code

When we try to write the static check code, visualize the sample code's SSA form at first.

We can use the below code to visualize the SSA form.

```golang
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
	if _, err := ssaPkg.WriteTo(os.Stdout); err != nil {
		log.Fatalf("Failed to write ssaPkg: cause %v", err)
	}
	for _, member := range ssaPkg.Members {
		if fn, ok := member.(*ssa.Function); ok {
			if fmt.Sprintf("./%s", fset.Position(fn.Pos()).Filename) != fset.Position(f.Pos()).Filename {
				continue
			}
			if _, err := ssaPkg.Func(fn.Name()).WriteTo(os.Stdout); err != nil {
				log.Fatalf("Failed to write func: cause %v", err)
			}
		}
	}
}
```

Let's try it to the red test code.

```golang
❯  go run ./utils/visualizessa -f ./testdata/src/a/a.go -d ./testdata/src/a
package a:
  type  T          struct{I int}
  func  funcB      func() (*T, error)
  func  init       func()
  var   init$guard bool

# Name: a.funcB
# Package: a
# Location: testdata/src/a/a.go:7:6
func funcB() (*T, error):
0:                                                                entry P:0 S:2
        ; var err error @ 8:6 is nil:error
        ; var err error @ 9:5 is nil:error
        t0 = nil:error != nil:error                                        bool
        ; *ast.BinaryExpr @ 9:5 is t0
        if t0 goto 1 else 2
1:                                                              if.then P:1 S:0
        return nil:*T, nil:error
2:                                                              if.done P:1 S:0
        return nil:*T, nil:error
```

(...after try and error) The first solution is here.

```
func run(pass *analysis.Pass) (any, error) {
	s := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	for _, f := range s.SrcFuncs {
		for _, b := range f.Blocks {
			for _, instr := range b.Instrs {
				ifstmt, ok := instr.(*ssa.BinOp)
				if !ok {
					continue
				}
				if isNilErrorCheck(ifstmt) {
					retBlock := ifstmt.Block().Succs[0]
					checkErrorReturnValue(retBlock, pass)
				}
			}
		}
	}
	return nil, nil
}

func isNilErrorCheck(ifstmt *ssa.BinOp) bool {
	if ifstmt.Op != token.NEQ {
		return false
	}
	return (isTypeError(ifstmt.Y.Type()) && ifstmt.X.Name() == nilErr) || (isTypeError(ifstmt.X.Type()) && ifstmt.Y.Name() == "nil:error")
}

func isTypeError(t types.Type) bool {
	if _, ok := t.Underlying().(*types.Interface); !ok {
		return false
	}
	return types.Identical(t, errorType)
}

func checkErrorReturnValue(b *ssa.BasicBlock, pass *analysis.Pass) {
	for _, instr := range b.Instrs {
		ret, ok := instr.(*ssa.Return)
		if !ok {
			continue
		}

		hasErr := false
		for _, v := range ret.Results {
			if isTypeError(v.Type()) && v.Name() != nilErr {
				hasErr = true
			}
		}
		if len(ret.Results) != 0 && !hasErr {
			pass.Reportf(ret.Pos(), "`return err` should be included in this return stmt. you seem to be ignoring error handling")
		}
	}
}
```

Test it.

```
❯ make test
go test
PASS
ok      github.com/MokkeMeguru/retnonnilerr     2.443s
```

Congratulation!

## More TDD

We got the analyzer to detect `return nil` despite of handling `err`.

But (you already know) this analyzer is not the completed one.

We need to more test cases and try TDD cycle to improve the degree of perfection.
