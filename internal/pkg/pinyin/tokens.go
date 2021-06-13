package pinyin

import "strings"

type Tokens []Token

func (t Tokens) String() string {
	var sb strings.Builder
	for i, tok := range t {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(tok.String())
	}
	return sb.String()
}
