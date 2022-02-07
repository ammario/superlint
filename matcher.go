package superlint

import (
	"os"
	"path/filepath"
)

// FileMatcherFunc describes a function that checks whether a file is relevant to
// a rule.
type FileMatcherFunc func(fi *os.File) bool

// ShellMatch uses standard shell syntax.
// E.g "**/*.go" matches all Go files.
func ShellMatch(pattern string) FileMatcherFunc {
	return func(fi *os.File) bool {
		m, err := filepath.Match(pattern, fi.Name())
		return m || err != nil
	}
}

// AndMatcher returns true when all match the file.
func AndMatcher(fm ...FileMatcherFunc) FileMatcherFunc {
	return func(fi *os.File) bool {
		for _, matcher := range fm {
			if !matcher(fi) {
				return false
			}
		}
		return true
	}
}

// OrMatcher returns true when any match the file.
func OrMatcher(fm ...FileMatcherFunc) FileMatcherFunc {
	return func(fi *os.File) bool {
		for _, matcher := range fm {
			if matcher(fi) {
				return true
			}
		}
		return false
	}
}
