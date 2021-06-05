package refdata

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/marcelocantos/flac/internal/fcache"
	"github.com/marcelocantos/flac/internal/pinyin"
	"github.com/marcelocantos/flac/internal/proto/cedict"
	"github.com/spf13/afero"
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

func loadCEDict(cache pinyin.Cache, fs afero.Fs, paths ...string) (*cedict.Dict, error) {
	result := &cedict.Dict{
		Simplified: map[string]*cedict.Entries{},
		Syllables:  map[string]bool{},
	}

	for _, path := range paths {
		target := &cedict.Dict{
			Simplified: map[string]*cedict.Entries{},
			Syllables:  map[string]bool{},
		}
		if err := fcache.Proto(fs, path, target, func(src io.Reader) error {
			data, err := afero.ReadAll(src)
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
					m := cedictDefRE.FindStringSubmatch(line)
					if len(m) != 5 {
						return lineError("no match")
					}
					// traditional := m[1]
					simplified := m[2]
					parts := strings.Split(m[3], " ")
					defs := m[4]

					for _, p := range parts {
						pinyin, err := cache.Pinyin(p)
						if err != nil {
							continue scanning
						}
						target.Syllables[pinyin.Syllable()] = true
					}

					entries, has := target.Simplified[simplified]
					if !has {
						entries = &cedict.Entries{Entries: map[string]*cedict.Entry{}}
						target.Simplified[simplified] = entries
					}

					answer := m[3]
					answer = strings.ReplaceAll(answer, " ", "")
					answer = strings.ReplaceAll(answer, "u:", "v")
					entry, has := entries.Entries[answer]
					if !has {
						entry = &cedict.Entry{}
						entries.Entries[answer] = entry
					}

					entry.Entry = append(entry.Entry, strings.Split(defs, "/")...)
				}
			}
			return nil
		}); err != nil {
			return nil, err
		}

		// Merge target into result.
		for word, entries := range target.Simplified {
			resultEntries, has := result.Simplified[word]
			if !has {
				resultEntries = &cedict.Entries{Entries: map[string]*cedict.Entry{}}
				result.Simplified[word] = resultEntries
			}
			for pinyin, entry := range entries.Entries {
				resultEntry, has := resultEntries.Entries[pinyin]
				if !has {
					resultEntry = &cedict.Entry{}
					resultEntries.Entries[pinyin] = resultEntry
				}
				resultEntry.Entry = append(resultEntry.Entry, entry.Entry...)
			}
		}
		for syllable, on := range target.Syllables {
			result.Syllables[syllable] = result.Syllables[syllable] || on
		}
	}
	return result, nil
}
