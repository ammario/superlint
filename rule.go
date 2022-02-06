package superlint

// A Rule defines a lint rule.
// The engine first assue
type Rule struct {
	Name        string
	FileMatcher FileMatcherFunc
	// Validator runs when all matched files have been loaded.
	// Validator runs all at once so that cross-codebase checks can be performed.
	// Validators should not return an error for lint violations.
	Validator ValidatorFunc
}

// RuleSet describes a set of rules loaded at linttime.
type RuleSet []Rule

// Add adds a new rule to the RuleSet.
func (rs *RuleSet) Add(r Rule) {
	*rs = append(*rs, r)
}
