package pinyin

import (
	"fmt"
	"regexp"
)

var (
	pinyinsRE = regexp.MustCompile(`(?i)^([a-zA-ZüÜ]+)([1-5]+)$`)
)

type Cache map[string]Pinyin

func (c Cache) MustPinyin(raw string) Pinyin {
	p, err := c.Pinyin(raw)
	if err != nil {
		panic(err)
	}
	return p
}

func (c Cache) Pinyin(raw string) (Pinyin, error) {
	p, has := c[raw]
	if !has {
		var err error
		p, err = newPinyin(raw)
		if err != nil {
			return Pinyin{}, err
		}
		c[raw] = p
	}
	return p, nil
}

func (c Cache) Pinyins(raw string) (Pinyins, error) {
	g := pinyinsRE.FindStringSubmatch(raw)
	if g == nil {
		return nil, fmt.Errorf("%s: not valid pinyin", raw)
	}
	result := make([]Pinyin, 0, len(g[2]))
	tones := 0
	for _, d := range g[2] {
		tones |= 1 << (d - '0')
	}
	for i := 1; i <= 5; i++ {
		if tones&(1<<i) != 0 {
			p, err := newPinyin(fmt.Sprintf("%s%d", g[1], i))
			if err != nil {
				return nil, err
			}
			result = append(result, p)
		}
	}
	return Pinyins(result), nil
}
