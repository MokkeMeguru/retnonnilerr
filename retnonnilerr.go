package retnonnilerr

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"github.com/gostaticanalysis/comment/passes/commentmap"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ssa"
)

const (
	name   = "retnonnilerr"
	doc    = "Retnonnilerr is a static analysis tool to detect `if err != nil { return nil }`"
	nilErr = "nil:error"
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
	samplefunc2()
	ignoredLines := getIgnoredLines(pass.Files, pass.Fset)
	s := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	for _, f := range s.SrcFuncs {
		fset := pass.Fset
		fileName := fset.Position(f.Pos()).Filename
		for _, b := range f.Blocks {
			for _, instr := range b.Instrs {
				ifstmt, ok := instr.(*ssa.BinOp)
				if !ok {
					continue
				}
				if isNilErrorCheck(ifstmt) {
					retBlock := ifstmt.Block().Succs[0]
					checkErrorReturnValue(retBlock, pass, fileName, ignoredLines)
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

func checkErrorReturnValue(b *ssa.BasicBlock, pass *analysis.Pass, filename string, ignoredLines map[string]map[int]bool) {
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
		if ignoredLines[filename][pass.Fset.Position(ret.Pos()).Line] {
			return
		}
		if len(ret.Results) != 0 && !hasErr {
			pass.Reportf(ret.Pos(), "`return err` should be included in this return stmt. you seems to throw the error handling")
		}
	}
}

func getIgnoredLines(files []*ast.File, fset *token.FileSet) map[string]map[int]bool {
	ignoredLines := make(map[string]map[int]bool)

	for _, file := range files {
		ignoredLinesByFile := make(map[int]bool)
		for _, group := range file.Comments {
			for _, comment := range group.List {
				if strings.Contains(comment.Text, "lint:ignore retnonnilerr") {
					position := fset.Position(comment.Pos())
					ignoredLinesByFile[position.Line+1] = true
				}
			}
		}
		fileName := fset.Position(file.Pos()).Filename
		ignoredLines[fileName] = ignoredLinesByFile
	}
	return ignoredLines
}

func samplefunc() error {
	return nil
}

func samplefunc2() error {
	if err := samplefunc(); err != nil {
		return nil
	}
	return nil
}
