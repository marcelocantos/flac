package outcome

import (
	"fmt"
	"sort"

	"github.com/marcelocantos/flac/internal/pkg/pinyin"
	"github.com/marcelocantos/flac/internal/pkg/proto/refdata"
)

type Outcome struct {
	Word       string
	Entries    *refdata.CEDict_Entries
	Bad        int
	Missing    int
	AnswerAlts pinyin.Alts
	Easy       bool
}

func (o *Outcome) Good() bool {
	return o.Bad == 0 && o.Missing == 0
}

func (o *Outcome) ErrorMessage() string {
	return fmt.Sprintf(
		"❌ %s ≠ %s\034❌ [#999999::]%[1]s ≠ [#999999::d]%[3]s[-::-]",
		o.Word, o.AnswerAlts.ColorString(), o.AnswerAlts.String())
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
