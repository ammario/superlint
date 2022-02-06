package superlint

import (
	"os"
)

// A Rule defines a lint rule.
// The engine first assue
type Rule struct {
	Name string
	// FileMatcher is written in filepath.Match syntax.
	// E.g "**/*.go" matches all Go files.
	FileMatcher string
	// Validator runs when all matched files have been loaded.
	Validator func(file []*os.File, info os.FileInfo) error
}
