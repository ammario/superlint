# superlint

`superlint` is an experimental, language-agnostic framework for lint rules written in Go.

[![Go Reference](https://pkg.go.dev/badge/github.com/ammario/superlint.svg)](https://pkg.go.dev/github.com/ammario/superlint)

superlint is designed to be a **superset of all possible linters**. Rules are
* codebase-scoped (as opposed to file or block scoped)
* Defined by arbitary Go code
  * Language-agonstic
  * Fast by default 
    * AST is parsed lazily
    * Linters run concurrently
  * Capable of using the network, filesystem, etc.
  * Composable

The vast ecosystem of existing linters can be called by a `superlint` ruleset. For example, a rule
can import an AST parser or execute a linting command.

`superlint` makes it easy to enforce arbitrary codebase-wide rules. For example, you may enforce that:
* Each Go binary has an accompanying Make entry
* Each http.Handler has an accompanying test
* Bash scripts don't exceed 1000 lines of code
* TypeScript/JS code is only in the `site` folder
* The `database` package never imports the `api` package

## Basic Usage

1. Create a rules file in your project (e.g `example/rules.go`)

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
    Linter: Single(lintgo.Validate(func(ps *lintgo.ParseState, _ FileInfo, report ReportFunc) error {
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
$ go install ./cmd/superlint && superlint -v example/rules.go
[18:05:36.552] loaded 2 rules
no-dog-files: example/dogs.go: no dogs allowed!
no-md5: example/dogs.go: crypto/md5 is insecure
        example/dogs.go:3       import "crypto/md5"
[18:05:36.560] 2 violations found
exit status 1

```

```

2. Run the rules

```bash
$ go build -buildmode=plugin -o rules.so example/rules.go && go run github.com/ammario/superlint/cmd/superlint rules.so
[18:05:36.552] loaded 2 rules
no-dog-files: example/dogs.go: no dogs allowed!
no-md5: example/dogs.go: crypto/md5 is insecure
        example/dogs.go:3       import "crypto/md5"
[18:05:36.560] 2 violations found
exit status 1

```



## Architecture

`superlint` loads your ruleset as a Go plugin. This is the only way `superlint` can support arbitrary Go lint rules
without direct integration with a build toolchain.