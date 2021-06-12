package assess

import (
	"fmt"
	"sort"
	"strings"

	"github.com/marcelocantos/flac/internal/pkg/pinyin"
	"github.com/marcelocantos/flac/internal/pkg/proto/refdata"
)

type Outcome struct {
	good       bool
	alts       pinyin.Alts
	correction string
}

func (o *Outcome) Alts() pinyin.Alts {
	return o.alts
}

func Assess(
	word string,
	entries *refdata.CEDict_Entries,
	answer string,
) *Outcome {
	o := &Outcome{good: true}

	tokenses, err := pinyin.Lex(answer)
	var answerAlts pinyin.Alts
	if err != nil {
		o.good = false
	} else {
		if len([]rune(word)) == 1 {
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
						o.good = false
						break
					}
					word = append(word, alts[0]...)
				}
				answerAlts = append(answerAlts, word)
			}
		}
		if o.good {
			o.alts = answerAlts

			altMap := map[string]bool{}
			covered := map[string]bool{}
			for _, alt := range answerAlts {
				altMap[alt.RawString()] = true
			}
			for raw := range entries.Definitions {
				if altMap[raw] {
					covered[raw] = true
					continue
				}
				lraw := strings.ToLower(raw)
				if altMap[lraw] {
					covered[lraw] = true
					continue
				}
				o.good = false
				break
			}
			if len(covered) != len(altMap) {
				o.good = false
			}
		}
	}

	if !o.good {
		// log.Printf("%v != %v", alts.RawString(), entries.Definitions)
		alts := make(pinyin.Alts, 0, len(entries.Definitions))
		for raw := range entries.Definitions {
			word, err := pinyin.NewWord(raw)
			if err != nil {
				panic(err)
			}
			alts = append(alts, word)
		}
		sort.Sort(alts)
		o.correction = fmt.Sprintf(
			// […l] normally means "blink", but we hijacked it for strikeout.
			// See cmd/flac/terminfo.go for details.
			"❌ %s ≠ %s (%[1]s = %[3]s)\034❌ [#999999::]%[1]s ≠ [#999999::d]%[4]s[-::-]",
			word, answerAlts.ColorString(), alts.ColorString(), answerAlts.String())
	}
	return o
}

func (o *Outcome) IsGood() bool {
	return o.good
}

func (o *Outcome) Correction() string {
	return o.correction
}
