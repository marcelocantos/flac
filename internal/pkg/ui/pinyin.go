package ui

import (
	"regexp"

	"github.com/gdamore/tcell/v2"
	"github.com/marcelocantos/flac/internal/pkg/pinyin"
	"github.com/rivo/tview"
)

var (
	inputRE = regexp.MustCompile(
		`(?i)^(?:[a-z]+[1-5]+(?:\s*(?:[/,·]\s*)?))*?(([a-z]+)[1-5]*|)$`)
	inputCharRE = regexp.MustCompile(`([a-z]+)[1-5]+`)
)

type PinyinInput struct {
	*tview.InputField

	syllables map[string]bool
	prefixes  map[string]bool

	submit func(answer string)
	giveUp func()
}

func newPinyinInput() *PinyinInput {
	input := &PinyinInput{
		InputField: tview.NewInputField(),
		syllables:  map[string]bool{},
		prefixes:   map[string]bool{},
		submit:     func(string) {},
		giveUp:     func() {},
	}
	input.SetAcceptanceFunc(input.accept)
	input.SetDoneFunc(input.done)
	return input
}

func (pi *PinyinInput) SetValidSyllables(syllables map[string]bool) *PinyinInput {
	for s := range syllables {
		// log.Println(s)
		pi.syllables[s] = true
		for i := 1; i <= len(s); i++ {
			pi.prefixes[s[:i]] = true
		}
	}
	// panic("")
	return pi
}

func (pi *PinyinInput) SetSubmit(submit func(answer string)) *PinyinInput {
	pi.submit = submit
	return pi
}

func (pi *PinyinInput) SetGiveUp(giveUp func()) *PinyinInput {
	pi.giveUp = giveUp
	return pi
}

func (pi *PinyinInput) accept(textToCheck string, lastChar rune) bool {
	m := inputRE.FindStringSubmatch(textToCheck)
	if m == nil {
		return false
	}
	for _, m := range inputCharRE.FindAllStringSubmatch(textToCheck, -1) {
		if _, err := (pinyin.Cache{}.WordAlts(m[0])); err != nil {
			return false
		}
		if !pi.syllables[m[1]] {
			return false
		}
	}
	return true
}

func (pi *PinyinInput) done(key tcell.Key) {
	switch key {
	case tcell.KeyEnter:
		text := pi.GetText()
		m := inputRE.FindStringSubmatch(text)
		if _, err := (pinyin.Cache{}).WordAlts(m[1]); err == nil {
			pi.submit(text)
		}
	case tcell.KeyEscape:
		pi.giveUp()
	}
}
