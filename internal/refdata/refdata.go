package refdata

import (
	"github.com/marcelocantos/flac/internal/pinyin"
	"github.com/marcelocantos/flac/internal/proto/cedict"
	"github.com/spf13/afero"
)

//go:generate curl https://www.plecoforums.com/download/global_wordfreq-release_utf-8-txt.2593/ | head -n 10000 | awk '//{print $1}' > words.txt
const (
	wordsFilename   = "words.txt"
	cedictFilename  = "cedict_1_0_ts_utf-8_mdbg.txt"
	addendaFilename = "addenda.txt"
)

type Refdata struct {
	words  []string
	cedict *cedict.Dict
}

func New(cache pinyin.Cache, fs afero.Fs) (*Refdata, error) {
	words, err := loadWords(fs, wordsFilename)
	if err != nil {
		return nil, err
	}

	cedict, err := loadCEDict(cache, fs, cedictFilename, addendaFilename)
	if err != nil {
		return nil, err
	}

	return &Refdata{
		words:  words,
		cedict: cedict,
	}, nil
}

func (d *Refdata) Words() []string {
	return d.words
}

func (d *Refdata) CEDict() *cedict.Dict {
	return d.cedict
}
