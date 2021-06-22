package ui

import (
	"github.com/rivo/tview"

	"github.com/marcelocantos/flac/internal/pkg/data"
	"github.com/marcelocantos/flac/internal/pkg/proto/refdata"
)

type Root struct {
	*tview.Flex

	Results *Results
	Answer  *AnswerInput
}

func New(db *data.Database, rd *refdata.RefData) *Root {
	results := newResults(db, rd)
	results.ScrollToEnd()

	input := newAnswerInput()

	keyboardHelp := " [orange::]❲?❳[-::]reveal [orange::]❲esc❳[-::]exit "
	keyboardHints := tview.NewTextView().
		SetDynamicColors(true).
		SetText(keyboardHelp)

	inputFlex := tview.NewFlex().
		AddItem(input, 0, 1, true).
		AddItem(keyboardHints, tview.TaggedStringWidth(keyboardHelp), 0, true)

	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(results, 0, 1, false).
		AddItem(inputFlex, 1, 0, true)

	return &Root{
		Flex:    mainFlex,
		Results: results,
		Answer:  input,
	}
}
