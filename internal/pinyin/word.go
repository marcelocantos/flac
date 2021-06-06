package pinyin

import (
	"strings"
)

type Word []Pinyin

func (w Word) Less(v Word) bool {
	n := len(w)
	if n > len(v) {
		n = len(v)
	}
	for i, a := range w[:n] {
		b := v[i]
		if a != b {
			return a.Less(b)
		}
	}
	return len(w) < len(v)
}

func (w Word) String() string {
	return w.string(func(p Pinyin) string { return p.String() })
}

func (w Word) ColorString() string {
	return w.string(func(p Pinyin) string { return p.ColorString() })
}

func (w Word) RawString() string {
	return w.string(func(p Pinyin) string { return p.RawString() })
}

func (w Word) string(s func(p Pinyin) string) string {
	var sb strings.Builder
	for i, p := range w {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(s(p))
	}
	return sb.String()
}
