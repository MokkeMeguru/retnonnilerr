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
├── cmd // cli で呼び出すツール
│   └── {package-name}
│       └── main.go
├── go.mod
├── {package-name}.go // 本体のコード
├── {package-name}_test.go
├── staticcheck.conf // ref: https://staticcheck.io/docs/configuration/
└── testdata // テストデータ
    └── src
        └── a
            ├── a.go
            └── a_test.go

```

### initial main code

the first code for static analyzer is following one.

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

var errorType = types.Universe.Lookup("error").Type()

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

	"github.com/MeguruMokke/retnonnilerr"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test_Run(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, retnonnilerr.Analyzer, "a")
}
```

i recommend you add the Makefile as following to test our analyzer easier.

```makefile
deps:
	go install github.com/kisielk/errcheck@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

lint:
	go vet ./...
	errcheck ./...
	staticcheck ./...

test:
	go test $(go list ./... | grep -v /testdata)
```

## add the first test data

let's write the first test code. (our analyzer should passin the test.)

```golang
package a

type T struct {
	I int
}

func funcA() (*T, error) {
	return nil, nil
}
```

test it.

```
❯ make test
go test
PASS
ok      github.com/MeguruMokke/retnonnilerr     1.349s
```

and then, add the red (failed) testcode.

```golang
func funcB() (*T, error) {
	var err error
	if err != nil {
		return nil, nil // want "return err"
	}
	return nil, nil
}
```

test it.

```
❯ make test
go test
--- FAIL: Test_Run (0.92s)
    analysistest.go:520: a/a.go:14: no diagnostic was reported matching `return err`
FAIL
exit status 1
FAIL    github.com/MeguruMokke/retnonnilerr     1.087s
make: *** [test] Error 1
```

we failed the tests. that means, we are in the TDD's (red-green-refactor) process now.

## write the test's green code

if we try to write the static check code, visualize the ast tree of sample code at first.

we can use the sample code to visualize the ast tree using the below utility code.

```golang
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
	ast.Print(fset, f)
}
```

let's try it to the test code.

```golang
❯ go run ./utils/visualizeast -f ./testdata/src/a/a.go
     0  *ast.File {
     1  .  Package: ./testdata/src/a/a.go:1:1
     2  .  Name: *ast.Ident {
     3  .  .  NamePos: ./testdata/src/a/a.go:1:9
     4  .  .  Name: "a"
     5  .  }
     6  .  Decls: []ast.Decl (len = 3) {
     7  .  .  0: *ast.GenDecl {
     ...
    53  .  .  }
    54  .  .  1: *ast.FuncDecl {
     ...
   112  .  .  }
   113  .  .  2: *ast.FuncDecl {
   114  .  .  .  Name: *ast.Ident {
   115  .  .  .  .  NamePos: ./testdata/src/a/a.go:11:6
   116  .  .  .  .  Name: "funcB"
   117  .  .  .  .  Obj: *ast.Object {
   118  .  .  .  .  .  Kind: func
   119  .  .  .  .  .  Name: "funcB"
   120  .  .  .  .  .  Decl: *(obj @ 113)
   121  .  .  .  .  }
   122  .  .  .  }
   123  .  .  .  Type: *ast.FuncType {
     ...
   151  .  .  .  }
   152  .  .  .  Body: *ast.BlockStmt {
   153  .  .  .  .  Lbrace: ./testdata/src/a/a.go:11:26
   154  .  .  .  .  List: []ast.Stmt (len = 3) {
   155  .  .  .  .  .  0: *ast.DeclStmt {
     ...
   182  .  .  .  .  .  }
   183  .  .  .  .  .  1: *ast.IfStmt {
   184  .  .  .  .  .  .  If: ./testdata/src/a/a.go:13:2
   185  .  .  .  .  .  .  Cond: *ast.BinaryExpr {
   186  .  .  .  .  .  .  .  X: *ast.Ident {
   187  .  .  .  .  .  .  .  .  NamePos: ./testdata/src/a/a.go:13:5
   188  .  .  .  .  .  .  .  .  Name: "err"
   189  .  .  .  .  .  .  .  .  Obj: *(obj @ 166)
   190  .  .  .  .  .  .  .  }
   191  .  .  .  .  .  .  .  OpPos: ./testdata/src/a/a.go:13:9
   192  .  .  .  .  .  .  .  Op: !=
   193  .  .  .  .  .  .  .  Y: *ast.Ident {
   194  .  .  .  .  .  .  .  .  NamePos: ./testdata/src/a/a.go:13:12
   195  .  .  .  .  .  .  .  .  Name: "nil"
   196  .  .  .  .  .  .  .  }
   197  .  .  .  .  .  .  }
   198  .  .  .  .  .  .  Body: *ast.BlockStmt {
   199  .  .  .  .  .  .  .  Lbrace: ./testdata/src/a/a.go:13:16
   200  .  .  .  .  .  .  .  List: []ast.Stmt (len = 1) {
   201  .  .  .  .  .  .  .  .  0: *ast.ReturnStmt {
   202  .  .  .  .  .  .  .  .  .  Return: ./testdata/src/a/a.go:14:3
   203  .  .  .  .  .  .  .  .  .  Results: []ast.Expr (len = 2) {
   204  .  .  .  .  .  .  .  .  .  .  0: *ast.Ident {
   205  .  .  .  .  .  .  .  .  .  .  .  NamePos: ./testdata/src/a/a.go:14:10
   206  .  .  .  .  .  .  .  .  .  .  .  Name: "nil"
   207  .  .  .  .  .  .  .  .  .  .  }
   208  .  .  .  .  .  .  .  .  .  .  1: *ast.Ident {
   209  .  .  .  .  .  .  .  .  .  .  .  NamePos: ./testdata/src/a/a.go:14:15
   210  .  .  .  .  .  .  .  .  .  .  .  Name: "nil"
   211  .  .  .  .  .  .  .  .  .  .  }
   212  .  .  .  .  .  .  .  .  .  }
   213  .  .  .  .  .  .  .  .  }
   214  .  .  .  .  .  .  .  }
   215  .  .  .  .  .  .  .  Rbrace: ./testdata/src/a/a.go:15:2
   216  .  .  .  .  .  .  }
   217  .  .  .  .  .  }
   218  .  .  .  .  .  2: *ast.ReturnStmt {
   219  .  .  .  .  .  .  Return: ./testdata/src/a/a.go:16:2
   220  .  .  .  .  .  .  Results: []ast.Expr (len = 2) {
   221  .  .  .  .  .  .  .  0: *ast.Ident {
   222  .  .  .  .  .  .  .  .  NamePos: ./testdata/src/a/a.go:16:9
   223  .  .  .  .  .  .  .  .  Name: "nil"
   224  .  .  .  .  .  .  .  }
   225  .  .  .  .  .  .  .  1: *ast.Ident {
   226  .  .  .  .  .  .  .  .  NamePos: ./testdata/src/a/a.go:16:14
   227  .  .  .  .  .  .  .  .  Name: "nil"
   228  .  .  .  .  .  .  .  }
   229  .  .  .  .  .  .  }
   230  .  .  .  .  .  }
   231  .  .  .  .  }
   232  .  .  .  .  Rbrace: ./testdata/src/a/a.go:17:1
   233  .  .  .  }
   234  .  .  }
   235  .  }
   ...
   256  }
```

the first solution of code is here.(some detail codes omitted)

```
func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
        // ignore test code
		if strings.Contains(f.Name.Name, "_test") {
			continue
		}
		for _, decl := range f.Decls {
			walkDecl(decl, false, pass)
		}
	}
	return nil, nil
}

func walkDecl(decl ast.Decl, hasNotNilError bool, pass *analysis.Pass) {
	switch decl := decl.(type) {
	case *ast.FuncDecl:
		walkStmt(decl.Body, hasNotNilError, pass)
	default:
		return
	}
}

func walkStmt(stmt ast.Stmt, hasNotNilError bool, pass *analysis.Pass) {
	switch stmt := stmt.(type) {
	case *ast.BlockStmt:
		for _, st := range stmt.List {
			walkStmt(st, hasNotNilError, pass)
		}
	case *ast.IfStmt:
		walkStmt(stmt.Body, isIfStmtValidateNilError(stmt, pass), pass)
	case *ast.ReturnStmt:
		if hasNotNilError {
			if !hasErrorInExprs(stmt.Results, pass) {
				pass.Reportf(stmt.Pos(), "`return err` should be included in this return stmt. you seems to throw the error handling")
			}
		}
	}
}

func isIfStmtValidateNilError(ifStmt *ast.IfStmt, pass *analysis.Pass) bool {
	switch cond := ifStmt.Cond.(type) {
	case *ast.BinaryExpr:
		switch cond.Op {
		case token.NEQ:
			return isErrorType(findExprType(cond.X, pass)) && isExprNil(cond.Y) ||
				isErrorType(findExprType(cond.Y, pass)) && isExprNil(cond.X)
		default:
			return false
		}
	default:
		return false
	}
}
```

let's test analyzer.

```
❯ make test
go test
PASS
ok      github.com/MeguruMokke/retnonnilerr     2.443s
```

congratulation!

## more TDD

we got the analyzer to detect `return nil` despite of handling `err`.

but (you already know) this analyzer is not the completed one.

we need to more test cases and try TDD cycle to improve the degree of perfection.

## include the code

use the analyzer as the linter, we install the analyzer code as command.

```golang
// cmd/retnonnilerr/main.go
package main

import (
	"github.com/MeguruMokke/retnonnilerr"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(retnonnilerr.Analyzer)
}
```

```
❯ go install github.com/MeguruMokke/retnonnilerr/cmd/retnonnilerr
❯ cd path/to/product
❯ retnonnilerr ./internal/...
```
