package ui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/marcelocantos/flac/internal/pkg/pinyin"
)

var (
	tradcharRE = regexp.MustCompile(`(?:\p{Han}+\|)(\p{Han}+)`)

	pinyinsRE = regexp.MustCompile(
		`(?i)(\p{Han}+)?\[((?:\w+\d\s+)*\w+\d)\]`)

	classifierRE = regexp.MustCompile(`\bCL:(\p{Han}+)`)

	classifierForRE = regexp.MustCompile(`\bclassifier for\b`)
)

func decorateDefinitions(defs []string) []string {
	ret := make([]string, 0, len(defs))

	toPrefix := "to "
	tos := make([]string, 0, len(defs))
	firstTo := -1

	for i, def := range defs {
		if strings.HasPrefix(def, toPrefix) {
			tos = append(tos, strings.TrimPrefix(def, toPrefix))
			if firstTo == -1 {
				ret = append(ret, "") // placeholder
				firstTo = i
			}
		} else {
			ret = append(ret, decorateDefinition(def))
		}
	}
	if firstTo != -1 {
		if len(tos) == 1 {
			ret[firstTo] = toPrefix + tos[0]
		} else {
			ret[firstTo] = fmt.Sprintf("%s[#888888::]⟨[-::]%s[#888888::]⟩[-::]",
				toPrefix, strings.Join(tos, "[#888888::],[-::] \035"))
		}
	}
	return ret
}

func decorateDefinition(phrase string) string {
	phrase = tradcharRE.ReplaceAllString(phrase, "$1")
	phrase = pinyinsRE.ReplaceAllStringFunc(phrase, func(s string) string {
		m := pinyinsRE.FindStringSubmatch(s)
		m[2] = pinyin.MustNewWord(m[2]).ColorString()
		return fmt.Sprintf("%s[%s]", m[1], m[2])
	})
	phrase = classifierRE.ReplaceAllString(phrase, "🆑:$1")
	phrase = classifierForRE.ReplaceAllString(phrase, "🆑 ➤")
	return phrase
}
