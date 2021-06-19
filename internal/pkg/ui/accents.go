package ui

import (
	"fmt"
	"regexp"

	"github.com/marcelocantos/flac/internal/pkg/pinyin"
)

var (
	tradcharRE = regexp.MustCompile(`(?:\p{Han}+\|)(\p{Han}+)`)

	pinyinsRE = regexp.MustCompile(
		`(?i)(\p{Han}+)?\[((?:\w+\d\s+)*\w+\d)\]`)
)

func accentPhrase(phrase string) string {
	phrase = tradcharRE.ReplaceAllString(phrase, "$1")
	phrase = pinyinsRE.ReplaceAllStringFunc(phrase, func(s string) string {
		m := pinyinsRE.FindStringSubmatch(s)
		m[2] = pinyin.MustNewWord(m[2]).ColorString()
		return fmt.Sprintf("%s[%s]", m[1], m[2])
	})
	return phrase
}
