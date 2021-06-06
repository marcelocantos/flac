package ui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"

	"github.com/marcelocantos/flac/internal/pkg/data"
)

type Root struct {
	*tview.Flex

	Results *Results
	Input   *PinyinInput
}

func New(db *data.Database) *Root {
	results := newResults(db)
	results.ScrollToEnd()
	fmt.Fprintf(results, "%s你好！", strings.Repeat("\n", 999))

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
