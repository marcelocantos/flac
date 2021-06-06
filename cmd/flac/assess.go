package main

import (
	"regexp"
	"sort"

	"github.com/marcelocantos/flac/internal/pinyin"
	"github.com/marcelocantos/flac/internal/proto/refdata"
)

var (
	alternativeSepRE = regexp.MustCompile(`\s*/\s*`)
)

type Outcome struct {
	good       bool
	pinyins    pinyin.Alts
	correction string
}

func Assess(pincache pinyin.Cache, entries *refdata.CEDict_Entries, answer string) *Outcome {
	o := &Outcome{good: true}

	words := alternativeSepRE.Split(answer, -1)
	alts := pinyin.Alts{}
	for _, word := range words {
		a, err := pincache.WordAlts(word)
		if err != nil {
			o.good = false
		}
		alts = append(alts, a...)
	}
	sort.Sort(alts)
	o.pinyins = alts

	if o.good {
		if len(alts) != len(entries.Definitions) {
			o.good = false
		} else {
			for _, p := range alts {
				if entries.Definitions[p.RawString()] == nil {
					o.good = false
					break
				}
			}
		}
	}
	if !o.good {
		// log.Printf("%v != %v", alts.RawString(), entries.Definitions)
		pinyins := make(pinyin.Alts, 0, len(entries.Definitions))
		for word := range entries.Definitions {
			pinyins = append(pinyins, pinyin.Word{pincache.MustNewPinyinNoResidue(word)})
		}
		sort.Sort(pinyins)
		o.correction = pinyins.ColorString()
	}
	return o
}

func (o *Outcome) IsGood() bool {
	return o.good
}

func (o *Outcome) Correction() string {
	return o.correction
}
