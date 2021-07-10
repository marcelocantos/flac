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
	root := &Root{}

	root.Results = newResults(db, rd)
	root.Results.ScrollToEnd()

	root.Hint = tview.NewTextView().
		SetDynamicColors(true)

	root.Answer = newAnswerInput()
	root.Answer.
		SetValidSyllables(rd.Dict.ValidSyllables).
		SetChangedFunc(func(text string) {
			if text != "" {
				root.Results.ClearMessages()
			}
		})

	keyboardHelp := " [orange::]❲?❳[-::]reveal [orange::]❲esc❳[-::]exit "
	keyboardHints := tview.NewTextView().
		SetDynamicColors(true).
		SetText(keyboardHelp)

	inputFlex := tview.NewFlex().
		AddItem(root.Answer, 0, 1, true).
		AddItem(keyboardHints, tview.TaggedStringWidth(keyboardHelp), 0, true)

	root.Flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(root.Results, 0, 1, false).
		AddItem(root.Hint, 1, 0, false).
		AddItem(inputFlex, 1, 0, true)

	return root
}
