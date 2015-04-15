package vcfgo

import (
	"fmt"
	"strings"
)

// VCFError satisfies the error interface and allows multiple errors.
// This is useful because, for example, on a single line, every sample may have
// a field that doesn't match the description in the header. We want to keep parsing
// but also let the caller know about the error.
type VCFError struct {
	Msgs  []string
	Lines []int64
}

// Error returns a string with all errors delimited by newlines.
func (e *VCFError) Error() string {
	var msgs []string
	seen := make(map[string]struct{})
	for i, m := range e.Msgs {
		// remove duplicates
		if _, ok := seen[m]; !ok {
			seen[m] = struct{}{}
			msgs = append(msgs, fmt.Sprintf("%s. [line: %d]", m, e.Lines[i]))
		}
	}
	return strings.Join(msgs, "\n")
}

// NewVCFError allocates the needed ingredients.
func NewVCFError() *VCFError {
	e := VCFError{Msgs: make([]string, 0), Lines: make([]int64, 0)}
	return &e
}

// Add adds an error and the line number within the vcf where the error took place.
func (e *VCFError) Add(err error, line int64) {
	if e != nil {
		if ierr := e.Error(); ierr != "" {
			e.Msgs = append(e.Msgs, ierr)
			e.Lines = append(e.Lines, line)
		}
	}
}

// IsEmpty returns true if there no errors stored.
func (e *VCFError) IsEmpty() bool {
	return len(e.Msgs) == 0
}

// Clear empties the Messages
func (e *VCFError) Clear() {
	e.Msgs = e.Msgs[:0]
	e.Lines = e.Lines[:0]
}
