package ui

import (
	"fmt"
)

type messages []string

func (m messages) write(format string, args ...interface{}) messages {
	return append(m, fmt.Sprintf(format, args...))
}
