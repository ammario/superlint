package superlint

import "github.com/coder/flog"

// RuleSet describes a set of rules loaded at linttime.
type RuleSet []Rule

// Add adds a new rule to the RuleSet.
func (rs *RuleSet) Add(r Rule) {
	*rs = append(*rs, r)
}

// Runner runs the RuleSet.
type Runner struct {
	DebugLogger *flog.Logger
	Log         *flog.Logger
}

func (rn *Runner) runRule(r Rule) {
}

func (rn *Runner) Run(matcher string, rs *RuleSet) error {
	for _, r := range *rs {
		rn.runRule(r)
	}
	return nil
}
