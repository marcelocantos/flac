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
}

func (o *Outcome) Good() bool {
	return o.Bad == 0 && o.Missing == 0
}

func (o *Outcome) Correction() string {
	wordAlts := make(pinyin.Alts, 0, len(o.Entries.Definitions))
	for raw := range o.Entries.Definitions {
		word, err := pinyin.NewWord(raw)
		if err != nil {
			panic(err)
		}
		wordAlts = append(wordAlts, word)
	}
	sort.Sort(wordAlts)

	return fmt.Sprintf(
		// […l] normally means "blink", but we hijacked it for strikeout.
		// See cmd/flac/terminfo.go for details.
		"❌ %s ≠ %s (%[1]s = %[3]s)\034❌ [#999999::]%[1]s ≠ [#999999::d]%[4]s[-::-]",
		o.Word, o.AnswerAlts.ColorString(), wordAlts.ColorString(), o.AnswerAlts.String())
}
