package pinyin

import "strings"

type Pinyins []Pinyin

func (p Pinyins) String() string {
	return p.string(func(p Pinyin) string { return p.String() })
}

func (p Pinyins) ColorString() string {
	return p.string(func(p Pinyin) string { return p.ColorString() })
}

func (p Pinyins) string(s func(p Pinyin) string) string {
	var sb strings.Builder
	for i, w := range p {
		if i > 0 {
			sb.WriteByte('/')
		}
		sb.WriteString(s(w))
	}
	return sb.String()
}

func (p Pinyins) Len() int {
	return len(p)
}

func (p Pinyins) Less(i, j int) bool {
	a, b := p[i], p[j]
	aLower := strings.ToLower(a.syllable)
	bLower := strings.ToLower(b.syllable)
	if aLower != bLower {
		return aLower < bLower
	}
	if a.syllable != b.syllable {
		return a.syllable < b.syllable
	}
	return a.tone < b.tone
}

func (p Pinyins) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
