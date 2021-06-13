package assess

import (
	"strings"

	"github.com/marcelocantos/flac/internal/pkg/outcome"
	"github.com/marcelocantos/flac/internal/pkg/pinyin"
	"github.com/marcelocantos/flac/internal/pkg/proto/refdata"
)

func Assess(
	word string,
	entries *refdata.CEDict_Entries,
	answer string,
) *outcome.Outcome {
	o := &outcome.Outcome{
		Word:    word,
		Entries: entries,
	}
	if answerAlts, ok := AnswerAlts(word, answer); ok {
		assess(entries, answerAlts, o)
	}
	return o
}

func AnswerAlts(word string, answer string) (pinyin.Alts, bool) {
	tokenses, err := pinyin.Lex(answer)
	if err != nil {
		return nil, false
	}
	var words []pinyin.Tokens
	var answerAlts pinyin.Alts
	if len([]rune(word)) == 1 {
		for _, tokens := range tokenses {
			words = append(words, tokens)
			for _, token := range tokens {
				answerAlts = append(answerAlts, token.Alts()...)
			}
		}
	} else {
		for _, tokens := range tokenses {
			var word pinyin.Word
			for _, token := range tokens {
				alts := token.Alts()
				if len(alts) != 1 {
					return nil, false
				}
				word = append(word, alts[0]...)
			}
			answerAlts = append(answerAlts, word)
		}
	}
	return answerAlts, true
}

func assess(entries *refdata.CEDict_Entries, answerAlts pinyin.Alts, o *outcome.Outcome) {
	answerMap := map[string]pinyin.Word{}
	for _, alt := range answerAlts {
		answerMap[alt.RawString()] = alt
	}

	defMap := map[string]bool{}
	for def := range entries.Definitions {
		defMap[strings.ToLower(def)] = true
	}

	partialDefs := map[string]bool{}

	for answer, alt := range answerMap {
		if defMap[answer] {
			o.Good = append(o.Good, alt)
		} else {
			tooShort := false
			for def := range defMap {
				if strings.HasPrefix(def, answer) {
					partialDefs[def] = true
					tooShort = true
				}
			}
			if tooShort {
				o.TooShort = append(o.TooShort, alt)
			} else {
				o.Bad = append(o.Bad, alt)
			}
		}
	}

	for def := range defMap {
		if _, has := answerMap[def]; !has {
			if !partialDefs[def] {
				o.Missing++
			}
		}
	}

	return
}
