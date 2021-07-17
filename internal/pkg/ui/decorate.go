package ui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/marcelocantos/flac/internal/pkg/pinyin"
)

var (
	taiwanRE = regexp.MustCompile(`Taiwan pr. `)

	tradcharRE = regexp.MustCompile(`(?:\p{Han}+\|)(\p{Han}+)`)

	pinyinsRE = regexp.MustCompile(
		`(?i)(\p{Han}+)?\[((?:\w+\d\s+)*\w+\d)\]`)

	classifierRE = regexp.MustCompile(`\bCL:(\p{Han}+)`)

	classifierForRE = regexp.MustCompile(`\bclassifier for\b`)
)

func decorateDefinitions(defs []string) []string {
	ret := make([]string, 0, len(defs))

	type group struct {
		prefix  string
		suffix  string
		replace string
		defs    []string
		first   int
	}
	groups := []*group{
		{first: -1, prefix: "to ", replace: "to… "},
		{first: -1, prefix: "abbr. for ", replace: "abbr… "},
		{first: -1, prefix: "classifier for "},
		{first: -1, prefix: "(grammatical equivalent of ", suffix: ")", replace: "(gramm ≣… "},
		{first: -1, prefix: "(indicates ", suffix: ")", replace: "(indic… "},
	}

defs:
	for _, def := range defs {
		for _, g := range groups {
			if strings.HasPrefix(def, g.prefix) && strings.HasSuffix(def, g.suffix) {
				g.defs = append(g.defs,
					strings.TrimSuffix(strings.TrimPrefix(def, g.prefix), g.suffix))
				if g.first == -1 {
					g.first = len(ret)
				} else {
					continue defs
				}
			}
		}
		ret = append(ret, def)
	}
	for _, g := range groups {
		if g.first != -1 {
			if len(g.defs) > 1 {
				ret[g.first] = fmt.Sprintf("%s\035%s%s",
					g.replace,
					strings.Join(g.defs, "[#666666::],[-::]\035"),
					g.suffix)
			}
		}
	}

	for i, def := range ret {
		ret[i] = DecorateDefinition(def)
	}

	return ret
}

func DecorateDefinition(phrase string) string {
	phrase = strings.ReplaceAll(phrase, "'", "’")
	phrase = taiwanRE.ReplaceAllString(phrase, "🇹🇼  ")
	phrase = tradcharRE.ReplaceAllString(phrase, "$1")
	phrase = pinyinsRE.ReplaceAllStringFunc(phrase, func(s string) string {
		m := pinyinsRE.FindStringSubmatch(s)
		m[2] = pinyin.MustNewWord(m[2]).ColorString("")
		return fmt.Sprintf("%s[%s]", m[1], m[2])
	})
	phrase = classifierRE.ReplaceAllString(phrase, "🆑:$1")
	phrase = classifierForRE.ReplaceAllString(phrase, "🆑➤")
	return phrase
}
