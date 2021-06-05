package pinyin

type Cache map[string]Pinyin

func (c Cache) Pinyin(raw string) (Pinyin, error) {
	p, has := c[raw]
	if !has {
		var err error
		p, err = NewPinyin(raw)
		if err != nil {
			return Pinyin{}, err
		}
		c[raw] = p
	}
	return p, nil
}
