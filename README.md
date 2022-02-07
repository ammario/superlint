# superlint

[![Go Reference](https://pkg.go.dev/badge/github.com/ammario/superlint.svg)](https://pkg.go.dev/github.com/ammario/superlint)

`superlint` is a linting system configured by user-defined Go code. Instead of a bespoke (and poorly documented) matching
language, `superlint` lets the user define arbitrary matching and validation functions.

`superlint` rules are **language-agonistic** and run **codebase-wide**. They're fast by default (lazy AST parsing) and
capable of  encorcing arbitrarily complex rules. A custom rule can create an AST for language-specific analysis.

For example, `superlint` can catch:
* Each Go binary has an accompanying Make entry
* Each http.Handler has an accompanying test
* Bash scripts don't exceed 1000 lines of code
* 

## Basic Usage

1. Create a rules file (e.g `superlint/rules.go`)

```go
package main

import (
	"os"
	"regexp"
	"strings"

	. "github.com/ammario/superlint"
	"github.com/ammario/superlint/lintgo"
	"github.com/coder/flog"
)

// LoadRules is the symbol loaded by superlint to inject rules.
var LoadRules Loader = func(_ *flog.Logger, r *RuleSet) {
	// `no-dog-files` rule uses no language awareness, and simply checks if `dog` exists in the filename.
	r.Add(Rule{
		Name: "no-dog-files",
		Validator: Single(func(fi *os.File, report ReportFunc) error {
			if strings.Contains(fi.Name(), "dog") {
				report(FileReference{}, "no dogs allowed!")
			}
			return nil
		}),
	})

	// `no-md5` demonstrates how language-aware features are possible in this paradigm.
	// lintgo is a simple wrapper around Go AST parsing.
	r.Add(Rule{
		Name:        "no-md5",
		FileMatcher: regexp.MustCompile(`\.go$`).MatchString,
		Validator: Single(lintgo.Validate(func(ps *lintgo.ParseState, _ *os.File, report ReportFunc) error {
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

```

2. Run the rules

```bash
I ~/P/a/superlint (master * =) go build -buildmode=plugin -o rules.so example/rules.go && go run github.com/ammario/superlint/cmd/superlint rules.so
[18:05:36.552] loaded 2 rules
no-dog-files: example/dogs.go: no dogs allowed!
no-md5: example/dogs.go: crypto/md5 is insecure
        example/dogs.go:3       import "crypto/md5"
[18:05:36.560] 2 violations found
exit status 1

```



## Architecture

`superlint` loads your ruleset as a Go plugin. This way `superlint` supports customizable lint rules without direct
integration with a build toolchain.