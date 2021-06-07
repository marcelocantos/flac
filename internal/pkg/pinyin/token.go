package pinyin

import (
	"strings"
)

type Token struct {
	syllable string
	tones    Tones
}

func (t Token) String() string {
	var sb strings.Builder
	sb.WriteString(t.syllable)
	sb.WriteByte('-')
	for r := t.tones.Range(); r.Next(); {
		sb.WriteByte(byte('0' + r.Value()))
	}
	return sb.String()
}

func (t Token) Alts() Alts {
	result := Alts{}
	for r := t.tones.Range(); r.Next(); {
		result = append(result,
			Word{newPinyinFromSyllableAndTone(t.syllable, r.Value())})
	}
	return result
}
