package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PinyinInput struct {
	*tview.InputField

	submit func(answer string)
	giveUp func()
}

func newPinyinInput() *PinyinInput {
	input := &PinyinInput{InputField: tview.NewInputField()}
	input.SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
		return true
	})
	input.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			input.submit(input.GetText())
		case tcell.KeyEscape:
			input.giveUp()
		}
	})
	return input
}

func (i *PinyinInput) SetSubmit(submit func(answer string)) {
	i.submit = submit
}

func (i *PinyinInput) SetGiveUp(giveUp func()) {
	i.giveUp = giveUp
}
