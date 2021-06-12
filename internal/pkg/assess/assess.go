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
	if err != nil {
		o.good = false
	} else {
		var alts pinyin.Alts
		if len([]rune(word)) == 1 {
			for _, tokens := range tokenses {
				for _, token := range tokens {
					alts = append(alts, token.Alts()...)
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
				alts = append(alts, word)
			}
		}
		if o.good {
			o.alts = alts

			altMap := map[string]bool{}
			covered := map[string]bool{}
			for _, alt := range alts {
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
		// \u200b = zero width space
		o.correction = fmt.Sprintf("%s\u200b = %s", word, alts.ColorString())
	}
	return o
}

func (o *Outcome) IsGood() bool {
	return o.good
}

func (o *Outcome) Correction() string {
	return o.correction
}
