package ui

import (
	"fmt"
	"regexp"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/marcelocantos/flac/internal/pkg/pinyin"
	"github.com/rivo/tview"
)

var (
	inputRE = regexp.MustCompile(
		`(?i)^(?:[a-z]+[1-5]+(?:\s*(?:[/,·]\s*)?))*?(([a-z]+)([1-5]*)|)\s*$`)
	inputCharRE = regexp.MustCompile(`([a-z]+)[1-5]+`)
)

type AnswerInput struct {
	*tview.InputField

	App *tview.Application

	compound bool

	syllables map[string]bool
	prefixes  map[string]bool

	submitFunc func(answer string)
	giveUpFunc func()
}

func newPinyinInput() *AnswerInput {
	input := &AnswerInput{
		InputField: tview.NewInputField(),
		syllables:  map[string]bool{},
		prefixes:   map[string]bool{},
		submitFunc: func(string) {},
		giveUpFunc: func() {},
	}
	input.SetAcceptanceFunc(input.accept)
	input.SetDoneFunc(input.done)
	return input
}

func (pi *AnswerInput) SetWord(word string, score int) {
	pi.SetLabel(fmt.Sprintf("%s[#999900::]%s[-::] ", word, brailleScore(score)))
	pi.compound = len([]rune(word)) > 1
}

func (pi *AnswerInput) SetValidSyllables(syllables []string) *AnswerInput {
	for _, s := range syllables {
		// log.Println(s)
		pi.syllables[s] = true
		for i := 1; i <= len(s); i++ {
			pi.prefixes[s[:i]] = true
		}
	}
	return pi
}

func (pi *AnswerInput) SetSubmitFunc(submit func(answer string)) *AnswerInput {
	pi.submitFunc = submit
	return pi
}

func (pi *AnswerInput) SetGiveUpFunc(giveUp func()) *AnswerInput {
	pi.giveUpFunc = giveUp
	return pi
}

func (pi *AnswerInput) FlashBackground() {
	pi.SetFieldBackgroundColor(tcell.ColorRed)
	go func() {
		time.Sleep(50 * time.Millisecond)
		pi.App.QueueUpdateDraw(func() {
			pi.SetFieldBackgroundColor(tview.Styles.ContrastBackgroundColor)
		})
	}()
}

func (pi *AnswerInput) accept(textToCheck string, lastChar rune) (ok bool) {
	defer func() {
		if !ok {
			pi.FlashBackground()
		}
	}()
	m := inputRE.FindStringSubmatch(textToCheck)
	if m == nil {
		return false
	}
	if m[2] != "" && !pi.prefixes[m[2]] {
		return false
	}
	for _, m := range inputCharRE.FindAllStringSubmatch(textToCheck, -1) {
		if !pi.syllables[m[1]] {
			return false
		}
	}
	return true
}

func (pi *AnswerInput) done(key tcell.Key) {
	switch key {
	case tcell.KeyEnter:
		text := pi.GetText()
		m := inputRE.FindStringSubmatch(text)
		if _, err := pinyin.WordAlts(m[1]); err == nil {
			pi.submitFunc(text)
		} else {
			pi.FlashBackground()
		}
	case tcell.KeyEscape:
		pi.giveUpFunc()
	}
}
