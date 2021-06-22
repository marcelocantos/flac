package main

import (
	"log"
	"regexp"

	"github.com/go-errors/errors"
	"github.com/marcelocantos/flac/internal/pkg/proto/refdata"
)

var (
	nonHanRE = regexp.MustCompile(`\P{Han}+`)

	wordSep = regexp.MustCompile(`[\s,，]+`)
)

func parsePhrase(phrase string, rd *refdata.RefData) ([]string, error) {
	phrase = nonHanRE.ReplaceAllString(phrase, "")

	var lengths []map[string]bool
	addWord := func(word string) {
		n := len(word)
		for len(lengths) < n+1 {
			lengths = append(lengths, map[string]bool{})
		}
		lengths[n][word] = true
	}
	for _, word := range rd.WordList.Words {
		addWord(word)
	}
	for word := range rd.Dict.Entries {
		addWord(word)
	}

	var ret []string
parsing:
	for len(phrase) > 0 {
		for length := len(lengths) - 1; length >= 1; length-- {
			if length == 1 {
				log.Print(lengths[length])
			}
			if length > len(phrase) {
				continue
			}
			candidate := phrase[:length]
			if lengths[length][candidate] {
				ret = append(ret, candidate)
				phrase = phrase[length:]
				continue parsing
			}
		}
		return nil, errors.Errorf("Unparsable content: %s", phrase)
	}

	return ret, nil
}

func parseWords(words string) []string {
	words = nonHanRE.ReplaceAllString(words, "")
	return wordSep.Split(words, -1)
}
