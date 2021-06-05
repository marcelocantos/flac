package ui

import (
	"github.com/rivo/tview"
)

type Results struct {
	*tview.TextView
}

func newResults() *Results {
	view := tview.NewTextView()
	view.SetDynamicColors(true)
	view.SetBorder(true)
	view.SetTitle("flac: learn 中文")

	return &Results{view}
}
