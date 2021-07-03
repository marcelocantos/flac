package main

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-errors/errors"
	"github.com/marcelocantos/flac/internal/pkg/pinyin"
	"github.com/marcelocantos/flac/internal/pkg/proto/refdata_pb"
	"github.com/spf13/afero"
)

func applyVariantRE(variantRE *regexp.Regexp, line string) (string, bool) {
	line2 := variantRE.ReplaceAllString(line, "$1")
	if line == line2 {
		return line, true
	}
	if backrefRE.MatchString(variantRE.String()) {
		groups := variantRE.FindStringSubmatch(line)
		for i, name := range variantRE.SubexpNames() {
			if name != "" {
				if j, err := strconv.ParseInt(name, 10, 64); err == nil {
					if groups[i] != groups[j] {
						return line, true
					}
				}
			}
		}
	}
	if strings.HasSuffix(line2, "] /") {
		return "", false
	}
	return line2, true
}

func applyVariantREs(line string) (string, bool) {
	interest := strings.Contains(line, `old variant of 和[he2]`)
	if interest {
		log.Print(line)
	}
	for _, variant := range []*regexp.Regexp{
		tradOnlyVariantRE,
		oldVariantRE,
		elidableVariantRE,
	} {
		var ok bool
		if line, ok = applyVariantRE(variant, line); !ok {
			if interest {
				log.Print(line, false)
			}
			return "", false
		}
	}
	if interest {
		log.Print(line, true)
	}
	return line, true
}

func loadCEDict(
	fs afero.Fs,
	path string,
	wm map[string]wordEntry,
	cedict *refdata_pb.CEDict,
) error {
	file, err := fs.Open(path)
	if err != nil {
		return err
	}

	syllables := cedict.ValidSyllables

	lineno := 0
	scanner := bufio.NewScanner(file)
scanning:
	for scanner.Scan() {
		lineno++
		if line := scanner.Text(); line != "" && !strings.HasPrefix(line, "#") {
			lineError := func(err error) error {
				return errors.WrapPrefix(err, fmt.Sprintf("%d: %s", lineno, line), 0)
			}

			var ok bool
			if line, ok = applyVariantREs(line); !ok {
				continue scanning
			}

			if match := cedictRemovalRE.FindStringSubmatch(line); match != nil {
				log.Print(match)
				simplified := match[2]
				if entries, has := cedict.Entries[simplified]; has {
					word, err := pinyin.NewWord(match[3])
					if err != nil {
						return lineError(err)
					}
					answer := word.RawString()
					if _, has := entries.Entries[answer]; has {
						log.Printf("  Removing %s", answer)
						delete(entries.Entries, answer)
						if len(entries.Entries) == 0 {
							log.Printf("  Removing %s", simplified)
							delete(cedict.Entries, simplified)
						}
					} else {
						log.Printf("  %s not found", answer)
					}
				} else {
					log.Printf("  %s not found", simplified)
				}
				continue
			}

			match := cedictDefRE.FindStringSubmatch(line)
			if match == nil {
				return lineError(errors.Errorf("no match"))
			}
			traditional := match[1]
			simplified := match[2]

			word, err := pinyin.NewWord(match[3])
			if err != nil {
				// log.Print(errors.WrapPrefix(err, fmt.Sprintf("%d: %s", lineno, match[3]), 0))
				continue
			}
			defs := match[4]

			for _, p := range word {
				syllables[p.Syllable()] = true
			}

			cedict.TraditionalToSimplified[traditional] = simplified

			entries, has := cedict.Entries[simplified]
			if !has {
				entries = &refdata_pb.CEDict_Entries{
					Entries: map[string]*refdata_pb.CEDict_Definitions{},
				}
				cedict.Entries[simplified] = entries
			}
			entries.Traditional = traditional

			answer := word.RawString()
			entry, has := entries.Entries[answer]
			if !has {
				entry = &refdata_pb.CEDict_Definitions{}
				entries.Entries[answer] = entry
			}

			entry.Definitions = append(entry.Definitions, strings.Split(defs, "/")...)
			simps, has := cedict.PinyinToSimplified[word.RawString()]
			if !has {
				simps = &refdata_pb.CEDict_Words{}
				cedict.PinyinToSimplified[word.RawString()] = simps
			}
			simps.Words = append(simps.Words, simplified)
		}
	}
	// log.Printf("Valid syllables: %v", cedict.ValidSyllables)

	maxReverse := 0
	var longestPinyinToSimplifiedPinyin string
	var longestPinyinToSimplified []string
	for word, words := range cedict.PinyinToSimplified {
		if maxReverse < len(words.Words) {
			maxReverse = len(words.Words)
			longestPinyinToSimplifiedPinyin = word
			longestPinyinToSimplified = words.Words
		}
	}
	log.Printf("Longest reverse mapping: %v 👉 (%d) %v",
		longestPinyinToSimplifiedPinyin,
		len(longestPinyinToSimplified),
		longestPinyinToSimplified)

	// Move "CL:..." and "...classifier for..." to the end.
	classifierForRE := regexp.MustCompile(`\bclassifier for `)
	for _, pred := range []func(string) bool{
		func(def string) bool { return classifierForRE.MatchString(def) },
		func(def string) bool { return strings.Contains(def, "CL:") },
	} {
		for _, entry := range cedict.Entries {
			for _, defs := range entry.Entries {
				for i := len(defs.Definitions) - 1; i >= 0; i-- {
					def := defs.Definitions[i]
					if pred(def) {
						defs.Definitions = append(defs.Definitions[:i], defs.Definitions[i+1:]...)
						defs.Definitions = append(defs.Definitions, def)
					}
				}
			}
		}
	}

	return scanner.Err()
}
