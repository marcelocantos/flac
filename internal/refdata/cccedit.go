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
	simplified map[string]*CeDictEntry
	syllables  map[string]bool
}

type CeDictEntry struct {
	Pinyin      []pinyin.Pinyin
	Definitions []string
}

func loadCeDict(fs afero.Fs, paths ...string) (*CEDict, error) {
	cedict := CEDict{
		simplified: map[string]*CeDictEntry{},
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

				pinyins := make([]pinyin.Pinyin, 0, len(parts))
				for _, p := range parts {
					pinyin, syllable, err := pinyin.NewPinyin(p)
					if err != nil {
						continue scanning
					}
					pinyins = append(pinyins, pinyin)
					cedict.syllables[syllable] = true
				}
				entry, has := cedict.simplified[simplified]
				if !has {
					entry = &CeDictEntry{Pinyin: pinyins}
					cedict.simplified[simplified] = entry
				}
				entry.Definitions = append(entry.Definitions, strings.Split(defs, "/")...)
			}
		}
	}
	return &cedict, nil
}
