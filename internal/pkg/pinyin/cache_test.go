package pinyin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWord(t *testing.T) {
	word, err := NewWord("wo3 men5")
	require.NoError(t, err)
	assert.Equal(t,
		Word{
			{pinyin: "wǒ", syllable: "wo", tone: 3},
			{pinyin: "men", syllable: "men", tone: 5},
		},
		word)
}
