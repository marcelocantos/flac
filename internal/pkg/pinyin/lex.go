package pinyin

import (
	"github.com/go-errors/errors"
)

func Lex(raw string) ([][]Token, error) {
	result := [][]Token{}
	tokens := []Token{}
	for raw != "" {
		groups := pinyinRE.FindStringSubmatch(raw)
		if groups == nil {
			return nil, errors.Errorf("%q: invalid pinyin", raw)
		}
		switch groups[1] {
		case "":
			tokens = append(tokens, Token{
				syllable: groups[2],
				tones:    newTonesFromString(groups[3]),
			})
		case "/":
			result = append(result, tokens)
			tokens = []Token{}
		default:
			tokens = append(tokens, Token{syllable: groups[1]})
		}
		raw = raw[len(groups[0]):]
	}
	result = append(result, tokens)
	return result, nil
}
