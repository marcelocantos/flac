package ui

import (
	"github.com/rivo/tview"

	"github.com/marcelocantos/flac/internal/pkg/data"
	"github.com/marcelocantos/flac/internal/pkg/proto/refdata_pb"
)

type Root struct {
	*tview.Flex

	Results *Results
	Hint    *tview.TextView
	Answer  *AnswerInput
}

func New(db *data.Database, rd *refdata_pb.RefData) *Root {
	results := newResults(db, rd)
	results.ScrollToEnd()

	hint := tview.NewTextView().
		SetDynamicColors(true)

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
		AddItem(hint, 1, 0, false).
		AddItem(inputFlex, 1, 0, true)

	return &Root{
		Flex:    mainFlex,
		Results: results,
		Hint:    hint,
		Answer:  input,
	}
}
