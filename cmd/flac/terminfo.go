package main

import (
	"os"

	"github.com/gdamore/tcell/v2/terminfo"
)

// Strikethrough is unavailable, so we hijack blink for the purpose.
func hackTerminfo() error {
	ti, err := terminfo.LookupTerminfo(os.Getenv("TERM"))
	if err != nil {
		return err
	}

	ti.Dim = "\033[9m"
	terminfo.AddTerminfo(ti)

	return nil
}
