package ui

import (
	"github.com/rivo/tview"
)

type Root struct {
	*tview.Flex

	Results *Results
	Input   *PinyinInput
}

func New() *Root {
	results := newResults()

	input := newPinyinInput()
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(results, 0, 1, false).
		AddItem(input, 1, 0, true)

	return &Root{
		Flex:    flex,
		Results: results,
		Input:   input,
	}
}
