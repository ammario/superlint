package superlint

import (
	"bytes"
	"fmt"
	"io"
)

// prettyPrintReference prints a human-readable description of the lint violation.
// It prints each line related to the violation and underlines the specific reference.
func prettyPrintReference(w io.Writer, fi []byte, reference FileReference) {
	// If the reference is beyond the length of the file, we have nothing to do.
	if reference.Pos.Offset >= len(fi) || reference.Pos.Offset >= len(fi) {
		return
	}
	if reference.Pos == reference.End {
		return
	}

	var (
		lines        [][]byte
		linesToWrite []int
		lastLineAt   int
	)
	for i, b := range fi {
		if b == '\n' {
			lines = append(lines, fi[lastLineAt:i])
			lastLineAt = i
		}
		// Write the current line
		switch {
		case i == reference.Pos.Offset:
			linesToWrite = append(linesToWrite, len(lines))
		case i > reference.Pos.Offset && i < reference.End.Offset && linesToWrite[len(linesToWrite)-1] != len(lines):
			linesToWrite = append(linesToWrite, len(lines))
		}
	}
	for _, ls := range linesToWrite {
		fmt.Fprintf(w, "\t%v:%v\t%s\n", reference.Name, ls, bytes.Trim(lines[ls], "\n"))
	}
}
