package pinyin

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-errors/errors"
)

var pinyinRE = regexp.MustCompile(`(?i)^\s*(?:([/,·])|([a-z:]+)([1-5]+))\s*`)

type Tone int8

func (t Tone) String() string {
	return strconv.Itoa(int(t))
}

var (
	toneColors = map[Tone]string{
		1: "red",
		2: "green",
		3: "blue",
		4: "purple",
		5: "black",
	}

	vowels = map[rune][]rune{ //  Breve forms:
		'a': []rune(" āáǎàa"), // ă
		'e': []rune(" ēéěèe"), // ĕ
		'i': []rune(" īíǐìi"), // ĭ
		'o': []rune(" ōóǒòo"), // ŏ
		'u': []rune(" ūúǔùu"), // ŭ
		'ü': []rune(" ǖǘǚǜü"), //
		'A': []rune(" ĀÁǍÀA"), // Ă
		'E': []rune(" ĒÉĚÈE"), // Ĕ
		'I': []rune(" ĪÍǏÌI"), // Ĭ
		'O': []rune(" ŌÓǑÒO"), // Ŏ
		'U': []rune(" ŪÚǓÙU"), // Ŭ
		'Ü': []rune(" ǕǗǙǛÜ"), //
	}
)

type Pinyin struct {
	pinyin   string
	syllable string
	tone     Tone
}

func NewPinyin(raw string) (_ Pinyin, residue string, err error) {
	groups := pinyinRE.FindStringSubmatch(raw)
	if groups == nil ||
		groups[1] == "" && len(groups[3]) != 1 {
		return Pinyin{}, "", errors.Errorf("%q: invalid pinyin", raw)
	}

	syllable := groups[2]
	syllable = strings.ReplaceAll(syllable, "v", "ü")
	syllable = strings.ReplaceAll(syllable, "u:", "ü")

	tone, err := strconv.Atoi(groups[3])
	if err != nil {
		return Pinyin{}, "", err
	}

	return newPinyinFromSyllableAndTone(syllable, tone), raw[len(groups[0]):], nil
}

func newPinyinFromSyllableAndTone(syllable string, tone int) Pinyin {
	chars := []rune(syllable)

	v := 0
	for _, c := range chars {
		if vowels[c] != nil {
			v++
		}
	}

	// https://en.wikipedia.org/wiki/Pinyin#Rules_for_placing_the_tone_mark
	switch {
	case v == 1:
		for i, c := range chars {
			if vowels[c] != nil {
				chars[i] = vowels[c][tone]
				break
			}
		}
	case strings.ContainsAny(syllable, "ae"):
		for i, c := range chars {
			if c == 'a' || c == 'e' {
				chars[i] = vowels[c][tone]
				break
			}
		}
	case strings.Contains(syllable, "ou"):
		for i, c := range chars {
			if c == 'o' {
				chars[i] = vowels[c][tone]
				break
			}
		}
	default:
		for i := len(chars) - 1; i >= 0; i-- {
			vowel := vowels[chars[i]]
			if vowel != nil {
				chars[i] = vowel[tone]
				break
			}
		}
	}
	return Pinyin{
		pinyin:   string(chars),
		syllable: strings.Replace(syllable, "ü", "v", 1),
		tone:     Tone(tone),
	}
}

func (p Pinyin) Less(q Pinyin) bool {
	aLower := strings.ToLower(p.syllable)
	bLower := strings.ToLower(q.syllable)
	if aLower != bLower {
		return aLower < bLower
	}
	if p.syllable != q.syllable {
		return p.syllable < q.syllable
	}
	return p.tone < q.tone
}

func (p Pinyin) String() string {
	return p.pinyin
}

func (p Pinyin) Syllable() string {
	return p.syllable
}

func (p Pinyin) Tone() Tone {
	return p.tone
}

func (p Pinyin) Color() string {
	return toneColors[p.Tone()]
}

func (p Pinyin) ColorString(flags string) string {
	return fmt.Sprintf("[%s::b%s]%s[-::-]", p.Color(), flags, p)
}

func (p Pinyin) RawString() string {
	return fmt.Sprintf("%s%d", p.syllable, p.tone)
}
