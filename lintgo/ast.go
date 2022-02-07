package lintgo

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"

	"github.com/ammario/superlint"
)

type ParseState struct {
	Fset *token.FileSet
	File *ast.File
}

// Validate traverses each files AST in depth-first order.
func Validate(
	fn func(ps *ParseState, fi *os.File, reporter superlint.ReportFunc) error,
) superlint.SingleValidator {
	fset := token.NewFileSet()
	return func(fi *os.File, report superlint.ReportFunc) error {
		goFile, err := parser.ParseFile(fset, fi.Name(), nil, parser.ParseComments)
		if err != nil {
			return err
		}
		return fn(&ParseState{fset, goFile}, fi, report)
	}
}
