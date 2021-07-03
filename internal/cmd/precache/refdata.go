package main

import (
	"io"
	"os"
	"regexp"
	"sort"

	"github.com/pierrec/lz4"
	"github.com/spf13/afero"
	"google.golang.org/protobuf/proto"

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
