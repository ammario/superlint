package main

import (
	"go/ast"
	"os"
	"strings"

	. "github.com/ammario/superlint"
	"github.com/ammario/superlint/lintgo"
	"github.com/coder/flog"
)

var LoadRules Loader = func(_ *flog.Logger, r *RuleSet) {
	r.Add(Rule{
		Name:        "no-dog-files",
		FileMatcher: ShellMatch("example/*.go"),
		Validator: Single(func(fi *os.File, report ReportFunc) error {
			if strings.Contains(fi.Name(), "dog") {
				report(FileReference{}, "no dogs allowed!")
			}
			return nil
		}),
	})

	r.Add(Rule{
		Name:        "no-zebra-functions",
		FileMatcher: ShellMatch("example/*.go"),
		Validator: Single(lintgo.Validate(func(goFile *ast.File, _ *os.File, report ReportFunc) error {
			ast.Inspect(goFile, func(node ast.Node) bool {
				funcCall, ok := node.(*ast.FuncDecl)
				if !ok {
					return true
				}
				funcName := funcCall.Name.Name
				if strings.Contains(funcName, "zebra") {
					report(
						FileReference{Pos: int(node.Pos()), End: int(node.End())},
						"no zebra allowed in a function name",
					)
				}
				return true
			})
			return nil
		}),
		),
	})
}
