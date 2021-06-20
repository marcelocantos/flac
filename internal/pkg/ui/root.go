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

	input := newPinyinInput()

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(results, 0, 1, false).
		AddItem(input, 1, 0, true)

	return &Root{
		Flex:    flex,
		Results: results,
		Answer:  input,
	}
}
