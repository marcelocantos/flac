package pinyin

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
