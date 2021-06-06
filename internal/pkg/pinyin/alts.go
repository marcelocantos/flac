package pinyin

import "strings"

type Alts []Word

func (a Alts) String() string {
	return a.string(func(w Word) string { return w.String() })
}

func (a Alts) ColorString() string {
	return a.string(func(w Word) string { return w.ColorString() })
}

func (a Alts) RawString() string {
	return a.string(func(w Word) string { return w.RawString() })
}

func (a Alts) string(s func(Word) string) string {
	var sb strings.Builder
	for i, w := range a {
		if i > 0 {
			sb.WriteString("/")
		}
		sb.WriteString(s(w))
	}
	return sb.String()
}

func (a Alts) Len() int {
	return len(a)
}

func (a Alts) Less(i, j int) bool {
	return a[i].Less(a[j])
}

func (a Alts) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
