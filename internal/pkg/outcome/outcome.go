package outcome

import (
	"fmt"
	"sort"

	"github.com/marcelocantos/flac/internal/pkg/pinyin"
	"github.com/marcelocantos/flac/internal/pkg/proto/refdata"
)

type Outcome struct {
	Word     string
	Entries  *refdata.CEDict_Entries
	Good     pinyin.Alts
	Bad      pinyin.Alts
	TooShort pinyin.Alts
	Missing  int
	Easy     bool
}

func (o *Outcome) Pass() bool {
	return len(o.Bad)+len(o.TooShort)+o.Missing == 0
}

func (o *Outcome) Correction() string {
	return fmt.Sprintf("%s = %s", o.Word, o.WordAlts().ColorString())
}

func (o *Outcome) WordAlts() pinyin.Alts {
	wordAlts := make(pinyin.Alts, 0, len(o.Entries.Definitions))
	for raw := range o.Entries.Definitions {
		word, err := pinyin.NewWord(raw)
		if err != nil {
			panic(err)
		}
		wordAlts = append(wordAlts, word)
	}
	sort.Sort(wordAlts)
	return wordAlts
}
