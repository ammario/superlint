package main

import (
	"os"
	"strings"

	. "github.com/ammario/superlint"
	"github.com/ammario/superlint/lintgo"
	"github.com/coder/flog"
)

// LoadRules is the symbol loaded by superlint to inject rules.
var LoadRules Loader = func(_ *flog.Logger, r *RuleSet) {
	// `no-dog-files` checks if `dog` exists in the filename.
	r.Add(Rule{
		Name: "no-dog-files",
		// "Single" here means that the rule does not need codebase-wide state.
		// Omit "Single" to receive all matching files.
		Linter: Single(func(fi FileInfo, report ReportFunc) error {
			if strings.Contains(fi.Name(), "dog") {
				report(FileReference{}, "no dogs allowed!")
			}
			return nil
		}),
	})

	// `no-md5` shows how language-awareness is possible in this paradigm.
	r.Add(Rule{
		Name: "no-md5",
		// lintgo is a simple wrapper around Go AST parsing.
		Linter: Single(lintgo.Validate(func(ps *lintgo.ParseState, _ os.FileInfo, report ReportFunc) error {
			for _, spec := range ps.File.Imports {
				if spec.Path.Value == "\"crypto/md5\"" {
					report(FileReference{
						Pos: ps.Fset.Position(spec.Path.Pos()).Offset,
						End: ps.Fset.Position(spec.Path.End()).Offset,
					}, "crypto/md5 is insecure")
				}
			}
			return nil
		}),
		),
	})
}
