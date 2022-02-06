package main

import (
	. "github.com/ammario/superlint"
	"github.com/coder/flog"
)

var LoadRules Loader = func(_ *flog.Logger, r *RuleSet) {
	r.Add(Rule{
		Name:        "no-http-in-database",
		FileMatcher: "../*.go",
	})
}
