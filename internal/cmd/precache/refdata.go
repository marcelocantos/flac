package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
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
	"github.com/marcelocantos/flac/internal/pkg/proto/refdata_pb"
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
		`^((?:\p{Han}+) (\p{Han}+) \[(.*?)\] /)` +
			`[^/]*(?:[^/]*\bvariant of|also written|see(?: also)?) ` +
			`(?:\p{Han}+\|)?(?P<2>\p{Han}+)\[(?P<3>.*?)\][^/]*/`)

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
	result := &refdata_pb.RefData{
		WordList: &refdata_pb.WordList{
			Frequencies: map[string]int64{},
		},
		Dict: &refdata_pb.CEDict{
			Entries:                 map[string]*refdata_pb.CEDict_Entries{},
			TraditionalToSimplified: map[string]string{},
			PinyinToSimplified:      map[string]*refdata_pb.CEDict_Words{},
		},
	}

	wordEntryMap, err := loadWords(fs, wordsPath, 10000)
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
	// if len(traditional) > 0 {
	// 	log.Printf("Eliding traditional words: %s", strings.Join(traditional, "  "))
	// }
	// if len(missing) > 0 {
	// 	log.Printf("Eliding words with no definitions: %s", strings.Join(missing, "  "))
	// }
	for _, word := range append(traditional, missing...) {
		delete(wordEntryMap, word)
	}

	words := make(wordEntries, 0, len(wordEntryMap))
	for _, entry := range wordEntryMap {
		words = append(words, entry)
	}
	sort.Sort(words)

	processWords(words, result.WordList)

	var writer io.WriteCloser
	writer, err = fs.Create(dest)
	if err != nil {
		return err
	}
	defer writer.Close()

	data, err := proto.Marshal(result)
	_ = data
	if err != nil {
		return err
	}
	if os.Getenv("FLAC_NO_COMPRESS") == "" {
		writer = lz4.NewWriter(writer)
		defer writer.Close()
	} else {
		writer.Write([]byte("NOCOMPRESS:"))
	}
	if _, err = writer.Write(data); err != nil {
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

func loadWords(fs afero.Fs, path string, limit int) (map[string]wordEntry, error) {
	wordsFile, err := fs.Open(path)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(bom.NewReader(wordsFile))
	words := map[string]wordEntry{}
	for i := 0; scanner.Scan() && (limit == -1 || i < limit); i++ {
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

func processWords(entries []wordEntry, wl *refdata_pb.WordList) {
	for _, entry := range entries {
		wl.Words = append(wl.Words, entry.word)
		wl.Frequencies[entry.word] = int64(entry.frequency)
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

	for s := range syllables {
		cedict.ValidSyllables = append(cedict.ValidSyllables, s)
	}
	sort.Strings(cedict.ValidSyllables)
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
