package main

import (
	"sort"
	"strings"

	"github.com/marcelocantos/flac/internal/pinyin"
	"github.com/marcelocantos/flac/internal/proto/refdata"
)

type Outcome struct {
	good       bool
	pinyins    pinyin.Pinyins
	correction string
}

func Assess(pincache pinyin.Cache, entries *refdata.CEDict_Entries, answer string) *Outcome {
	o := &Outcome{good: true}

	words := strings.Split(answer, "/")
	pinyins := pinyin.Pinyins{}
	for _, word := range words {
		p, err := pincache.Pinyins(word)
		if err != nil {
			o.good = false
		}
		pinyins = append(pinyins, p...)
	}
	sort.Sort(pinyins)
	o.pinyins = pinyins

	if o.good {
		if len(pinyins) != len(entries.Definitions) {
			o.good = false
		} else {
			for _, p := range pinyins {
				if entries.Definitions[p.RawString()] == nil {
					o.good = false
					break
				}
			}
		}
	}
	if !o.good {
		pinyins := make(pinyin.Pinyins, 0, len(entries.Definitions))
		for word := range entries.Definitions {
			pinyins = append(pinyins, pincache.MustPinyin(word))
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
