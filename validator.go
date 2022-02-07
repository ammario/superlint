package superlint

import (
	"fmt"
)

type ValidatorFunc func(files map[string]FileInfo, report ReportFunc) error

// SingleValidator a ValidatorFunc for rules that do not perform cross-file checks.
// SingleValidator populates FileReference.Name in the report automatically.
type SingleValidator func(f FileInfo, report ReportFunc) error

// Single forms a new SingleValidator.
func Single(vf SingleValidator) ValidatorFunc {
	return func(files map[string]FileInfo, report ReportFunc) error {
		for _, file := range files {
			err := vf(file, func(reference FileReference, message string) {
				reference.Name = file.Path
				report(reference, message)
			})
			if err != nil {
				return fmt.Errorf("%v: %w", err)
			}
		}
		return nil
	}
}
