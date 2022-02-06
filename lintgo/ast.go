package lintgo

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"

	"github.com/ammario/superlint"
)

// Validate traverses each files AST in depth-first order.
func Validate(
	fn func(fset *token.FileSet, goFile *ast.File, fi *os.File, reporter superlint.ReportFunc) error,
) superlint.SingleValidator {
	fset := token.NewFileSet()
	return func(fi *os.File, report superlint.ReportFunc) error {
		goFile, err := parser.ParseFile(fset, fi.Name(), nil, parser.ParseComments)
		if err != nil {
			return err
		}
		return fn(fset, goFile, fi, report)
	}
}
