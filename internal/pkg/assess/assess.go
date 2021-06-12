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
	if answerAlts, ok := answerAlts(len([]rune(word)) == 1, answer); ok {
		o.AnswerAlts = answerAlts
		o.Bad, o.Missing = assess(entries, answerAlts)
	}
	return o
}

func answerAlts(simple bool, answer string) (pinyin.Alts, bool) {
	tokenses, err := pinyin.Lex(answer)
	if err != nil {
		return nil, false
	}
	var answerAlts pinyin.Alts
	if simple {
		for _, tokens := range tokenses {
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

func assess(entries *refdata.CEDict_Entries, answerAlts pinyin.Alts) (bad, missing int) {
	answerMap := map[string]bool{}
	for _, alt := range answerAlts {
		answerMap[alt.RawString()] = true
	}

	defMap := map[string]bool{}
	for def := range entries.Definitions {
		defMap[strings.ToLower(def)] = true
	}

	for answer := range answerMap {
		if !defMap[answer] {
			bad++
		}
	}

	for def := range defMap {
		if !answerMap[def] {
			missing++
		}
	}

	return
}
