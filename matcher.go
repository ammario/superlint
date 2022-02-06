package superlint

import (
	"path/filepath"
)

// FileMatcherFunc describes a function that checks whether a file is relevant to
// a rule.
type FileMatcherFunc func(name string) bool

// ShellMatch uses standard shell syntax.
// E.g "**/*.go" matches all Go files.
func ShellMatch(pattern string) FileMatcherFunc {
	return func(name string) bool {
		m, err := filepath.Match(pattern, name)
		return m || err != nil
	}
}

// AndMatcher returns true when all match the file.
func AndMatcher(fm ...FileMatcherFunc) FileMatcherFunc {
	return func(name string) bool {
		for _, matcher := range fm {
			if !matcher(name) {
				return false
			}
		}
		return true
	}
}

// OrMatcher returns true when any match the file.
func OrMatcher(fm ...FileMatcherFunc) FileMatcherFunc {
	return func(name string) bool {
		for _, matcher := range fm {
			if matcher(name) {
				return true
			}
		}
		return false
	}
}
