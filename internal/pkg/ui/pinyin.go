package ui

import (
	"regexp"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	inputRE = regexp.MustCompile(
		`(?i)^[a-zü]+([1-5]+((?:/|\s+)?[a-zü]+[1-5]+)*((?:/|\s+)?([a-zü]+[1-5]?)?)?)?$`)
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
	return inputRE.Match([]byte(textToCheck))
	// d := len(textToCheck) - 1
	// if '0' <= textToCheck[d] && textToCheck[d] <= '5' {
	// 	i := strings.IndexAny(textToCheck[:d], "12345")
	// 	return pi.syllables[textToCheck[:i]]
	// }
	// return pi.prefixes[textToCheck]
}

func (pi *PinyinInput) done(key tcell.Key) {
	switch key {
	case tcell.KeyEnter:
		if pi.GetText() != "" {
			pi.submit(pi.GetText())
		}
	case tcell.KeyEscape:
		pi.giveUp()
	}
}
