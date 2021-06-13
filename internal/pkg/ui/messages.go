package ui

import (
	"fmt"
)

type messages []string

// BlankOutMessages sets all messages to the empty string rather than simply
// removing them. This clears them out without causing the view to scroll.
func (m messages) blankOut() (changed bool) {
	for i := range m {
		if len(m) > 0 {
			changed = true
		}
		m[i] = ""
	}
	return
}

func (m messages) write(format string, args ...interface{}) messages {
	return append(m, fmt.Sprintf(format, args...))
}
