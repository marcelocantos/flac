package refdata

import (
	"bufio"

	"github.com/spf13/afero"
)

func loadWords(fs afero.Fs, path string) ([]string, error) {
	wordsFile, err := fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer wordsFile.Close()

	scanner := bufio.NewScanner(wordsFile)
	var words []string
	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			words = append(words, line)
		}
	}
	return words, nil
}
