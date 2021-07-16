package ui

// import (
// 	"fmt"

// 	"github.com/gdamore/tcell/v2"
// 	"github.com/rivo/tview"
// )

// type DefinitionInput struct {
// 	*tview.InputField

// 	exitFunc   func()
// 	giveUpFunc func()
// 	submitFunc func(answer string)
// }

// func newDefinitionInput() *DefinitionInput {
// 	input := &DefinitionInput{
// 		InputField: tview.NewInputField(),
// 		exitFunc:   func() {},
// 		giveUpFunc: func() {},
// 		submitFunc: func(string) {},
// 	}
// 	input.SetDoneFunc(input.done)
// 	return input
// }

// func (pi *DefinitionInput) SetWord(word string, score int) int {
// 	label := fmt.Sprintf("%s[#999900::]%s[-::] ", word, brailleScore(score))
// 	pi.SetLabel(label)
// 	pi.compound = len([]rune(word)) > 1
// 	return tview.TaggedStringWidth(label)
// }

// func (pi *DefinitionInput) SetValidSyllables(syllables map[string]bool) *DefinitionInput {
// 	pi.syllables = syllables
// 	for s := range syllables {
// 		for i := 1; i <= len(s); i++ {
// 			pi.prefixes[s[:i]] = true
// 		}
// 	}
// 	return pi
// }

// func (pi *DefinitionInput) SetExitFunc(exitFunc func()) *DefinitionInput {
// 	pi.exitFunc = exitFunc
// 	return pi
// }

// func (pi *DefinitionInput) SetGiveUpFunc(giveUp func()) *DefinitionInput {
// 	pi.giveUpFunc = giveUp
// 	return pi
// }

// func (pi *DefinitionInput) SetSubmitFunc(submit func(answer string)) *DefinitionInput {
// 	pi.submitFunc = submit
// 	return pi
// }

// func (pi *DefinitionInput) done(key tcell.Key) {
// 	switch key {
// 	case tcell.KeyEnter:
// 		text := pi.GetText()
// 		pi.submitFunc(text)
// 	case tcell.KeyEscape:
// 		pi.exitFunc()
// 	}
// }
