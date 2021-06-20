package pinyin

import (
	"fmt"
	"strings"
)

type Alts []Word

func (a Alts) String() string {
	return a.string("", func(w Word) string { return w.String() })
}

func (a Alts) ColorString(flags string) string {
	return a.string(flags, func(w Word) string { return w.ColorString(flags) })
}

func (a Alts) RawString() string {
	return a.string("", func(w Word) string { return w.RawString() })
}

func (a Alts) string(flags string, s func(Word) string) string {
	var sb strings.Builder
	for i, w := range a {
		if i > 0 {
			if flags == "" {
				sb.WriteByte('/')
			} else {
				fmt.Fprintf(&sb, "[::%s]/[::-]", flags)
			}
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
