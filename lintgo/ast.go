package lintgo

import (
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"

	"github.com/ammario/superlint"
)

type ParseState struct {
	Fset *token.FileSet
	File *ast.File
}

var matchGoFiles = regexp.MustCompile(`\.go$`)

// Validate traverses each files AST in depth-first order.
func Validate(
	fn func(ps *ParseState, fi superlint.FileInfo, reporter superlint.ReportFunc) error,
) superlint.SingleValidator {
	fset := token.NewFileSet()
	return func(fi superlint.FileInfo, report superlint.ReportFunc) error {
		if fi.IsDir() || !matchGoFiles.MatchString(fi.Name()) {
			return nil
		}
		goFile, err := parser.ParseFile(fset, fi.Path, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		return fn(&ParseState{fset, goFile}, fi, report)
	}
}
