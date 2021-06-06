package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/pierrec/lz4"
	"github.com/spf13/afero"
	"google.golang.org/protobuf/proto"

	"github.com/marcelocantos/flac/internal/pinyin"
	"github.com/marcelocantos/flac/internal/proto/refdata"
	"github.com/marcelocantos/flac/internal/refdata/words"
)

var (
	cedictDefRE = regexp.MustCompile(
		`(\S+) (\S+) \[((?:[\w:]+ (?:(?:[\w:]+|[,·]) )*)?[\w:]+)\] /(.*)/$`)

	// Detect traditional-only variants.
	tradOnlyVariantRE = regexp.MustCompile(
		`^((?:.) (.) \[(.*?)\] /)(?:old )?variant of (?:.\|)?(?P<2>.)\[(?P<3>.*?)\]/`)

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
			Frequencies: map[string]int32{},
			Positions:   map[string]int32{},
		},
		Dict: &refdata.CEDict{
			Entries:                 map[string]*refdata.CEDict_Entries{},
			Syllables:               map[string]bool{},
			TraditionalToSimplified: map[string]string{},
		},
	}

	if err := loadWords(fs, wordsPath, result.WordList); err != nil {
		return err
	}

	wl := &words.WordList{WordList: result.WordList}

	for _, path := range dictPaths {
		if err := loadCEDict(fs, path, wl, result.Dict); err != nil {
			return err
		}
	}

	log.Printf("dest: %s", dest)
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

func loadWords(fs afero.Fs, path string, wl *refdata.WordList) error {
	words_data, err := afero.ReadFile(fs, path)
	if err != nil {
		return err
	}
	wordsFile := bytes.NewBuffer(words_data)

	scanner := bufio.NewScanner(wordsFile)
	i := -1
	for scanner.Scan() {
		i++
		if line := scanner.Text(); line != "" {
			parts := strings.SplitN(line, "\t", 2)
			word := parts[0]
			freq, err := strconv.Atoi(parts[1])
			if err != nil {
				return err
			}

			wl.Words = append(wl.Words, word)
			wl.Frequencies[word] = int32(freq)
			wl.Positions[word] = int32(i)
		}
	}
	return scanner.Err()
}

func loadCEDict(
	fs afero.Fs,
	path string,
	wl *words.WordList,
	cedict *refdata.CEDict,
) error {
	pincache := pinyin.Cache{}
	data, err := afero.ReadFile(fs, path)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(data))
	lineno := 0
scanning:
	for scanner.Scan() {
		lineno++
		if line := scanner.Text(); line != "" && !strings.HasPrefix(line, "#") {
			lineError := func(msg string) error {
				return fmt.Errorf("%s: %d: %s", msg, lineno, line)
			}

			match := cedictDefRE.FindStringSubmatch(line)
			if len(match) != 5 {
				return lineError("no match")
			}
			traditional := match[1]
			simplified := match[2]
			if wl.Has(simplified) {
				continue
			}

		variants:
			for _, variant := range []struct {
				re       *regexp.Regexp
				backrefs bool
			}{
				{tradOnlyVariantRE, true},
				{oldVariantRE, false},
				{elidableVariantRE, false},
			} {
				line2 := variant.re.ReplaceAllString(line, "$1")
				if line != line2 {
					if variant.backrefs {
						groups := variant.re.FindStringSubmatch(line)
						for i, name := range variant.re.SubexpNames() {
							if name != "" {
								if j, err := strconv.ParseInt(name, 10, 64); err == nil {
									if groups[i] != groups[j] {
										continue variants
									}
								}
							}
						}
					}
					if strings.HasSuffix(line2, "] /") {
						continue scanning
					}
					line = line2
				}
			}
			parts := strings.Split(match[3], " ")
			defs := match[4]

			for _, p := range parts {
				pinyin, err := pincache.Pinyin(p)
				if err != nil {
					continue scanning
				}
				cedict.Syllables[pinyin.Syllable()] = true
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

			answer := match[3]
			answer = strings.ReplaceAll(answer, " ", "")
			answer = strings.ReplaceAll(answer, "u:", "v")
			entry, has := entries.Definitions[answer]
			if !has {
				entry = &refdata.CEDict_Definitions{}
				entries.Definitions[answer] = entry
			}

			entry.Definitions = append(entry.Definitions, strings.Split(defs, "/")...)
		}
	}
	return scanner.Err()
}
