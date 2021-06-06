package pinyin

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-errors/errors"
)

var (
	pinyinRE = regexp.MustCompile(`(?i)^\s*(?:([,·])|([a-zü:]+)([1-5]+))\s*`)
)

type Tone int8

func (t Tone) String() string {
	return strconv.Itoa(int(t))
}

var (
	toneColors = map[Tone]string{
		1: "[red::b]",
		2: "[green::b]",
		3: "[blue::b]",
		4: "[purple::b]",
		5: "[black::b]",
	}

	vowels = map[rune][]rune{
		'a': []rune(" āáǎàa"),
		'e': []rune(" ēéěèe"),
		'i': []rune(" īíǐìi"),
		'o': []rune(" ōóǒòo"),
		'u': []rune(" ūúǔùu"),
		'ü': []rune(" ǖǘǚǜü"),
	}
)

type Pinyin struct {
	pinyin   string
	syllable string
	tone     Tone
}

func newPinyin(raw string) (_ Pinyin, residue string, _ error) {
	groups := pinyinRE.FindStringSubmatch(raw)
	if groups == nil ||
		groups[1] == "" && len(groups[3]) != 1 {
		return Pinyin{}, "", errors.Errorf("%q: invalid pinyin", raw)
	}
	syllable := groups[2]
	tone, err := strconv.Atoi(groups[3])
	if err != nil {
		panic(err)
	}
	syllable = strings.ReplaceAll(syllable, "v", "ü")
	syllable = strings.ReplaceAll(syllable, "u:", "ü")

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
	}, raw[len(groups[0]):], nil
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

func (p Pinyin) ColorString() string {
	return fmt.Sprintf("%s%s[-::-]", p.Color(), p)
}

func (p Pinyin) RawString() string {
	return fmt.Sprintf("%s%d", p.syllable, p.tone)
}
