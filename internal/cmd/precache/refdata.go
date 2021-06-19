package main

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/go-errors/errors"
	"github.com/pierrec/lz4"
	"github.com/spf13/afero"
	"github.com/spkg/bom"
	"google.golang.org/protobuf/proto"

	"github.com/marcelocantos/flac/internal/pkg/pinyin"
	"github.com/marcelocantos/flac/internal/pkg/proto/refdata"
)

var (
	backrefRE = regexp.MustCompile(`\(\?P<\d>`)

	hanziRE = regexp.MustCompile(`^\p{Han}+`)

	cedictDefRE = regexp.MustCompile(
		`(\S+) (\S+) \[((?:[\w:]+ (?:(?:[\w:]+|[,·]) )*)?[\w:]+)\] /(.*)/$`)

	cedictRemovalRE = regexp.MustCompile(
		`- (\S+) (\S+) \[((?:[\w:]+ (?:(?:[\w:]+|[,·]) )*)?[\w:]+)\]`)

	// Detect traditional-only variants.
	tradOnlyVariantRE = regexp.MustCompile(
		`^((?:.) (.) \[(.*?)\] /)[^/]*(?:old )?variant of (?:.\|)?(?P<2>.)\[(?P<3>.*?)\][^/]*/`)

	// Detect old variants.
	oldVariantRE = regexp.MustCompile(
		`(/)(?:\((?:old|archaic)\) [^/]*|[^/]* \((?:old|archaic)\)|(?:old|archaic) variant of [^/]*)/`)

	// Detect other elidable content.
	elidableVariantRE = regexp.MustCompile(
		`(/)[^/]*(?:\(dialect\)|Taiwan pr\.)[^/]*/`)
)

func cacheRefData(
	fs afero.Fs,
	wordsPath string,
	dictPaths []string,
	dest string,
) error {
	result := &refdata.RefData{
		WordList: &refdata.WordList{
			Frequencies: map[string]int64{},
			Positions:   map[string]int64{},
		},
		Dict: &refdata.CEDict{
			Entries:                 map[string]*refdata.CEDict_Entries{},
			TraditionalToSimplified: map[string]string{},
			PinyinToSimplified:      map[string]*refdata.CEDict_Words{},
		},
	}

	wordEntryMap, err := loadWords(fs, wordsPath)
	if err != nil {
		return err
	}

	for _, path := range dictPaths {
		if err := loadCEDict(fs, path, wordEntryMap, result.Dict); err != nil {
			return err
		}
	}

	traditional := []string{}
	missing := []string{}
	for _, entry := range wordEntryMap {
		word := entry.word
		if _, has := result.Dict.Entries[word]; !has {
			if _, has := result.Dict.TraditionalToSimplified[word]; has {
				traditional = append(traditional, word)
			} else {
				missing = append(missing, word)
			}
		}
	}
	if len(traditional) > 0 {
		log.Printf("Eliding traditional words: %s", strings.Join(traditional, "  "))
	}
	if len(missing) > 0 {
		log.Printf("Eliding words with no definitions: %s", strings.Join(missing, "  "))
	}
	for _, word := range append(traditional, missing...) {
		delete(wordEntryMap, word)
	}

	words := make(wordEntries, 0, len(wordEntryMap))
	for _, entry := range wordEntryMap {
		words = append(words, entry)
	}
	sort.Sort(words)

	processWords(words, result.WordList)

	out, err := fs.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	data, err := proto.Marshal(result)
	_ = data
	if err != nil {
		return err
	}
	w := lz4.NewWriter(out)
	defer w.Close()
	if _, err = w.Write(data); err != nil {
		return err
	}

	return nil
}

type wordEntry struct {
	word      string
	index     int
	frequency int
}

type wordEntries []wordEntry

func (e wordEntries) Len() int {
	return len(e)
}

func (e wordEntries) Less(i, j int) bool {
	a, b := e[i], e[j]
	return a.index < b.index
}

func (e wordEntries) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func loadWords(fs afero.Fs, path string) (map[string]wordEntry, error) {
	wordsFile, err := fs.Open(path)
	if err != nil {
		return nil, err
	}
	i := -1
	scanner := bufio.NewScanner(bom.NewReader(wordsFile))
	words := map[string]wordEntry{}
	for scanner.Scan() {
		i++
		if line := scanner.Text(); line != "" {
			parts := strings.SplitN(line, "\t", 2)
			word := parts[0]
			if !hanziRE.MatchString(word) {
				continue
			}
			freq, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, err
			}

			words[word] = wordEntry{word: word, index: i, frequency: freq}
		}
	}
	return words, scanner.Err()
}

func processWords(entries []wordEntry, wl *refdata.WordList) {
	for i, entry := range entries {
		wl.Words = append(wl.Words, entry.word)
		wl.Frequencies[entry.word] = int64(entry.frequency)
		wl.Positions[entry.word] = int64(i)
	}
}

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
	for _, variant := range []*regexp.Regexp{
		tradOnlyVariantRE,
		oldVariantRE,
		elidableVariantRE,
	} {
		var ok bool
		if line, ok = applyVariantRE(variant, line); !ok {
			return "", false
		}
	}
	return line, true
}

func loadCEDict(
	fs afero.Fs,
	path string,
	wm map[string]wordEntry,
	cedict *refdata.CEDict,
) error {
	file, err := fs.Open(path)
	if err != nil {
		return err
	}

	syllables := map[string]bool{}

	lineno := 0
	scanner := bufio.NewScanner(file)
scanning:
	for scanner.Scan() {
		lineno++
		if line := scanner.Text(); line != "" && !strings.HasPrefix(line, "#") {
			lineError := func(err error) error {
				return errors.WrapPrefix(err, fmt.Sprintf("%d: %s", lineno, line), 0)
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
					if _, has := entries.Definitions[answer]; has {
						log.Printf("  Removing %s", answer)
						delete(entries.Definitions, answer)
						if len(entries.Definitions) == 0 {
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
			if _, has := wm[simplified]; !has {
				continue
			}

			var ok bool
			if line, ok = applyVariantREs(line); !ok {
				continue scanning
			}

			word, err := pinyin.NewWord(match[3])
			if err != nil {
				return lineError(err)
			}
			defs := match[4]

			for _, p := range word {
				syllables[p.Syllable()] = true
			}

			cedict.TraditionalToSimplified[traditional] = simplified

			entries, has := cedict.Entries[simplified]
			if !has {
				entries = &refdata.CEDict_Entries{
					Definitions: map[string]*refdata.CEDict_Definitions{},
				}
				cedict.Entries[simplified] = entries
			}
			entries.Traditional = traditional

			answer := word.RawString()
			entry, has := entries.Definitions[answer]
			if !has {
				entry = &refdata.CEDict_Definitions{}
				entries.Definitions[answer] = entry
			}

			entry.Definitions = append(entry.Definitions, strings.Split(defs, "/")...)
			simps, has := cedict.PinyinToSimplified[word.RawString()]
			if !has {
				simps = &refdata.CEDict_Words{}
				cedict.PinyinToSimplified[word.RawString()] = simps
			}
			simps.Words = append(simps.Words, simplified)
		}
	}

	for s := range syllables {
		cedict.ValidSyllables = append(cedict.ValidSyllables, s)
	}
	sort.Strings(cedict.ValidSyllables)
	log.Printf("Valid syllables: %v", cedict.ValidSyllables)

	maxReverse := 0
	var longestPinyinToSimplified []string
	for _, words := range cedict.PinyinToSimplified {
		if maxReverse < len(words.Words) {
			maxReverse = len(words.Words)
			longestPinyinToSimplified = words.Words
		}
	}
	log.Printf("Longest reverse mapping: %v", longestPinyinToSimplified)

	return scanner.Err()
}
