package lintgo

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
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
	fn func(ps *ParseState, fi os.FileInfo, reporter superlint.ReportFunc) error,
) superlint.SingleValidator {
	fset := token.NewFileSet()
	return func(fi superlint.FileInfo, report superlint.ReportFunc) error {
		if fi.IsDir() || !matchGoFiles.MatchString(fi.Name()) {
			return nil
		}
		fmt.Printf("%v\n", fi.Name())
		goFile, err := parser.ParseFile(fset, fi.Name(), nil, parser.ParseComments)
		if err != nil {
			return err
		}
		return fn(&ParseState{fset, goFile}, fi, report)
	}
}
