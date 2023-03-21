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
	ignoredLines := getIgnoredLines(pass.Files, pass.Fset)
	s := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	for _, f := range s.SrcFuncs {
		fileName := pass.Fset.Position(f.Pos()).Filename
		for _, b := range f.Blocks {
			for _, instr := range b.Instrs {
				switch instr := instr.(type) {
				case *ssa.BinOp:
					if isNilErrorCheck(instr) {
						retBlock := instr.Block().Succs[0]
						checkErrorReturnValue(retBlock, pass, ignoredLines[fileName])
					}
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

func checkErrorReturnValue(b *ssa.BasicBlock, pass *analysis.Pass, ignoredLines map[int]bool) {
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
		if ignoredLines[pass.Fset.Position(ret.Pos()).Line] {
			return
		}
		if len(ret.Results) != 0 && !hasErr {
			pass.Reportf(ret.Pos(), "`return err` should be included in this return stmt. you seem to be ignoring error handling")
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
