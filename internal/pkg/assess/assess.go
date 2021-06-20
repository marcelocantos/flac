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
	var answerAlts pinyin.Alts
	if len([]rune(word)) == 1 {
		for _, tokens := range tokenses {
			for _, token := range tokens {
				answerAlts = append(answerAlts, token.Alts()...)
			}
		}
	} else {
		for _, tokens := range tokenses {
			altses := []pinyin.Alts{}
			for _, token := range tokens {
				altses = append(altses, token.Alts())
			}
			answerAlts = answerProduct(answerAlts, altses, pinyin.Word{})
		}
	}
	return answerAlts, true
}

func answerProduct(answerAlts pinyin.Alts, altses []pinyin.Alts, word pinyin.Word) pinyin.Alts {
	if len(altses) == 0 {
		return append(answerAlts, word)
	}
	for _, alt := range altses[0] {
		answerAlts = answerProduct(answerAlts, altses[1:], append(word, alt[0]))
	}
	return answerAlts
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
			badTones := false
			for def := range defMap {
				word := pinyin.MustNewWord(def)
				if len(alt) < len(word) && alt.RawString() == word[:len(alt)].RawString() {
					partialDefs[def] = true
					tooShort = true
				} else if len(alt) == len(word) {
					syllableErrors := 0
					tonalErrors := 0
					for i, p := range word {
						if alt[i].Syllable() != p.Syllable() {
							syllableErrors++
						}
						if alt[i].Tone() != p.Tone() {
							tonalErrors++
						}
					}
					if syllableErrors == 0 && tonalErrors > 0 {
						badTones = true
					}
				}
			}
			if tooShort {
				o.TooShort = append(o.TooShort, alt)
			} else if badTones {
				o.BadTones = append(o.BadTones, alt)
				o.Bad = append(o.Bad, alt)
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
}
