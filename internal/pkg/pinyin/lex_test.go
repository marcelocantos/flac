package pinyin

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLex(t *testing.T) {
	t.Parallel()

	assertLex(t, "", "")
	assertLex(t, "shi4", "shi-4")
	assertLex(t, "shi4de5", "shi-4 de-5")
	assertLex(t, "shi4de5", "shi-4 de-5")
	assertLex(t, "dou1/Du1/du1", "dou-1", "Du-1", "du-1")
	assertLex(t, "dou1/Du1/du1", "dou-1", "Du-1", "du-1")
	assertLex(t, "yi1 kong3 zhi1 jian4", "yi-1 kong-3 zhi-1 jian-4")
	assertLex(t, "xu1yao4/xiang3", "xu-1 yao-4", "xiang-3")
	assertLex(t, "xu1 yao4/xiang3", "xu-1 yao-4", "xiang-3")
	assertLex(t, " xu1 yao4 / xiang3 ", "xu-1 yao-4", "xiang-3")

	assertLex(t, "jiang14qiang1", "jiang-14 qiang-1")
	assertLex(t, "jiang14/qiang1", "jiang-14", "qiang-1")
}

func assertLex(t *testing.T, src string, expected ...string) bool {
	t.Helper()

	tokenses, err := Lex(src)
	var actual []string
	for _, tokens := range tokenses {
		var chars []string
		for _, token := range tokens {
			chars = append(chars, token.String())
		}
		actual = append(actual, strings.Join(chars, " "))
	}
	require.NoError(t, err)
	return assert.Equal(t, expected, actual)
}
