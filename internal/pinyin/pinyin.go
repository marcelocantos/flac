package pinyin

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	pinyinRE = regexp.MustCompile(`^([a-zü:]+)(\d)$`)
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
		4: "[magenta::b]",
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

	vowelTones = func() map[rune]Tone {
		result := make(map[rune]Tone, 6*5)
		for _, runes := range vowels {
			for i, r := range runes[1:] {
				result[r] = Tone(1 + i)
			}
		}
		return result
	}()
)

type Pinyin string

func NewPinyin(pinyin string) (_ Pinyin, syllable string, _ error) {
	if pinyin == "," {
		return Pinyin(", "), "", nil
	}
	if pinyin == "·" {
		return Pinyin(" · "), "", nil
	}
	groups := pinyinRE.FindStringSubmatch(pinyin)
	if groups == nil {
		return "", "", fmt.Errorf("%q not a valid pinyin form", pinyin)
	}
	syllable = groups[1]
	tone, err := strconv.Atoi(groups[2])
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
	return Pinyin(string(chars)), strings.Replace(syllable, "ü", "v", 1), nil
}

func (p Pinyin) Tone() Tone {
	for _, r := range p {
		if tone, has := vowelTones[r]; has {
			return Tone(tone)
		}
	}
	return 0
}

func (p Pinyin) Color() string {
	return toneColors[p.Tone()]
}

func (p Pinyin) ColorString() string {
	return fmt.Sprintf("%s%s[::]", p.Color(), p)
}
