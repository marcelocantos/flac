package refdata

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/marcelocantos/flac/internal/pinyin"
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

type CEDict struct {
	simplified map[string]CEDictEntries
	syllables  map[string]bool
}

func (d *CEDict) Simplified() map[string]CEDictEntries {
	return d.simplified
}

func (d *CEDict) Syllables() map[string]bool {
	return d.syllables
}

type CEDictEntries map[string]*CEDictEntry

type CEDictEntry []string

func loadCEDict(fs afero.Fs, paths ...string) (*CEDict, error) {
	cedict := CEDict{
		simplified: map[string]CEDictEntries{},
		syllables:  map[string]bool{},
	}

	for _, path := range paths {
		file, err := fs.Open(path)
		if err != nil {
			return nil, err
		}
		scanner := bufio.NewScanner(file)
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
					return nil, lineError("no match")
				}
				// traditional := m[1]
				simplified := m[2]
				parts := strings.Split(m[3], " ")
				defs := m[4]

				for _, p := range parts {
					pinyin, err := pinyin.NewPinyin(p)
					if err != nil {
						continue scanning
					}
					cedict.syllables[pinyin.Syllable()] = true
				}

				entries, has := cedict.simplified[simplified]
				if !has {
					entries = CEDictEntries{}
					cedict.simplified[simplified] = entries
				}

				answer := m[3]
				answer = strings.ReplaceAll(answer, " ", "")
				answer = strings.ReplaceAll(answer, "u:", "v")
				entry, has := entries[answer]
				if !has {
					entry = &CEDictEntry{}
					entries[answer] = entry
				}

				*entry = append(*entry, strings.Split(defs, "/")...)
			}
		}
	}
	return &cedict, nil
}
