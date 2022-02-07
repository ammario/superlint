package superlint

// FileReference describes a chunk of a file.
type FileReference struct {
	Name string
	// Pos, End are byte indices when specifying a problem within the file.
	Pos, End int
}

// ReportFunc describes a function used to report lint violations.
type ReportFunc func(reference FileReference, message string)
