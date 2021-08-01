package main

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/spf13/afero"
	"github.com/spkg/bom"

	"github.com/marcelocantos/flac/internal/pkg/proto/refdata_pb"
)

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
		wl.Frequencies[entry.word] = int32(entry.frequency)
	}
}
