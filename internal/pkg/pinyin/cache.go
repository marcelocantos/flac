package pinyin

import (
	"fmt"
	"regexp"

	"github.com/go-errors/errors"
)

var pinyinsRE = regexp.MustCompile(`(?i)^([a-zA-ZüÜ]+)([1-5]+)$`)

func MustNewPinyinNoResidue(raw string) Pinyin {
	p, residue, err := NewPinyin(raw)
	if err != nil {
		panic(err)
	}
	if residue != "" {
		panic(errors.Errorf("%q: invalid pinyin form", raw))
	}
	return p
}

func NewWord(raw string) (Word, error) {
	word := Word{}
	for residue := raw; residue != ""; {
		var p Pinyin
		var err error
		p, residue, err = NewPinyin(residue)
		if err != nil {
			return nil, errors.WrapPrefix(err, fmt.Sprintf("%q: invalid word form", raw), 0)
		}
		word = append(word, p)
	}
	return word, nil
}

func WordAlts(raw string) (Alts, error) {
	g := pinyinsRE.FindStringSubmatch(raw)
	if g == nil || len(g[0]) < len(raw) {
		return nil, errors.Errorf("%q: not valid pinyin", raw)
	}
	result := make(Alts, 0, len(g[2]))
	tones := 0
	for _, d := range g[2] {
		tones |= 1 << (d - '0')
	}
	for i := 1; i <= 5; i++ {
		if tones&(1<<i) != 0 {
			p, _, err := NewPinyin(fmt.Sprintf("%s%d", g[1], i))
			if err != nil {
				return nil, err
			}
			result = append(result, Word{p})
		}
	}
	return result, nil
}
