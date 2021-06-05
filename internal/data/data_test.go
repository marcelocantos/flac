package data_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/marcelocantos/flac/internal/data"
)

var rawwords = []string{
	"第", "的", "了", "在", "是", "我", "和", "有", "你", "个", "也", "这", "不",
	"他", "上", "人", "中", "就", "年", "为", "对", "说", "都", "要", "到", "着",
	"~", "与", "将", "日", "我们", "好", "月", "会", "大", "来", "还", "等", "而",
	"地", "自己", "后", "两", "-", "被", "没有", "去", "但", "从", "很", "给", "时",
	"以", "中国",
}

var words = []string{
	"第", "的", "了", "在", "是", "我", "和", "有", "你", "个", "也", "这", "不",
	"他", "上", "人", "中", "就", "年", "为", "对", "说", "都", "要", "到", "着",
	"与", "将", "日", "我们", "好", "月", "会", "大", "来", "还", "等", "而",
	"地", "自己", "后", "两", "被", "没有", "去", "但", "从", "很", "给", "时",
	"以", "中国",
}

func TestDatabasePopulate(t *testing.T) {
	t.Parallel()

	d, err := data.NewDatabase(":memory:")
	require.NoError(t, err)
	for n := 1; n < len(words)*2; n *= 2 {
		n := n
		if n > len(words) {
			n = len(words)
		}
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			slice := words[:n]
			require.NoError(t, d.Populate(slice))

			max, err := d.MaxPos()
			require.NoError(t, err)
			assert.Equal(t, n-1, max)

			for i, word := range slice {
				pos, err := d.WordPos(word)
				require.NoError(t, err)
				assert.Equal(t, i, pos)
			}
		})
	}
}

func TestWordAt(t *testing.T) {
	t.Parallel()
	d := prepareDatabase(t)

	word, err := d.WordAt(0)
	require.NoError(t, err)
	assert.Equal(t, "第", word)

	word, err = d.HeadWord()
	require.NoError(t, err)
	assert.Equal(t, "第", word)

	word, err = d.WordAt(1)
	require.NoError(t, err)
	assert.Equal(t, "的", word)
}

func TestDatabaseWordPos(t *testing.T) {
	t.Parallel()
	d := prepareDatabase(t)

	pos, err := d.WordPos("第")
	_ = assert.NoError(t, err) && assert.Equal(t, 0, pos)

	_, err = d.WordPos("元")
	assert.Error(t, err)
}

func TestDatabaseMoveWord(t *testing.T) {
	t.Parallel()
	d := prepareDatabase(t)

	require.NoError(t, d.MoveWord("第", 1))
	assertWordsAt(t, d, 0, "的", "第", "了")

	require.NoError(t, d.MoveWord("的", 3))
	assertWordsAt(t, d, 0, "第", "了", "在", "的")

	require.NoError(t, d.MoveWord("了", 3))
	assertWordsAt(t, d, 0, "第", "在", "的", "了")

	require.NoError(t, d.MoveWord("了", 1))
	assertWordsAt(t, d, 0, "第", "了", "在", "的")

	require.NoError(t, d.MoveWord("的", 0))
	assertWordsAt(t, d, 0, "的", "第", "了", "在")

	require.NoError(t, d.MoveWord("第", 0))
	assertWordsAt(t, d, 0, "第", "的", "了", "在")
}

func TestDatabaseMoveWordFromEnd(t *testing.T) {
	t.Parallel()
	d := prepareDatabase(t)

	require.NoError(t, d.MoveWord("中国", 10))
	assertWordsAt(t, d, 0, words[:10]...)
	assertWordsAt(t, d, 10, "中国")
	assertWordsAt(t, d, 11, words[10:len(words)-1]...)
}

func TestDatabaseMoveWordToEnd(t *testing.T) {
	t.Parallel()
	d := prepareDatabase(t)

	require.NoError(t, d.MoveWord("也", len(words)-1))
	assertWordsAt(t, d, 0, words[:10]...)
	assertWordsAt(t, d, 10, words[11:len(words)-1]...)
	assertWordsAt(t, d, len(words)-1, "也")
}

func TestDatabaseMoveWordPastEnd(t *testing.T) {
	t.Parallel()
	d := prepareDatabase(t)

	require.NoError(t, d.MoveWord("也", len(words)+99))
	assertWordsAt(t, d, 0, words[:10]...)
	assertWordsAt(t, d, 10, words[11:len(words)-1]...)
	assertWordsAt(t, d, len(words)-1, "也")
}

func TestDatabaseMoveWordFromStartToEnd(t *testing.T) {
	t.Parallel()
	d := prepareDatabase(t)

	require.NoError(t, d.MoveWord("第", len(words)-1))
	assertWordsAt(t, d, 0, words[1:len(words)-1]...)
	assertWordsAt(t, d, len(words)-1, "第")
}

func TestDatabaseMoveWordFromEndToStart(t *testing.T) {
	t.Parallel()
	d := prepareDatabase(t)

	require.NoError(t, d.MoveWord("中国", 0))
	assertWordsAt(t, d, 0, "中国")
	assertWordsAt(t, d, 1, words[:len(words)-1]...)
}

func TestDatabaseMoveMissingWord(t *testing.T) {
	t.Parallel()
	d := prepareDatabase(t)

	require.Error(t, d.MoveWord("元", 10))
}

func prepareDatabase(t *testing.T) *data.Database {
	d, err := data.NewDatabase(":memory:")
	require.NoError(t, err)
	require.NoError(t, d.Populate(rawwords))
	return d
}

func assertWordsAt(t *testing.T, d *data.Database, expected int, words ...string) bool {
	t.Helper()

	ok := true
	for i, word := range words {
		pos, err := d.WordPos(word)
		ok = ok &&
			assert.NoError(t, err) &&
			assert.Equal(t, expected+i, pos)
	}
	return ok
}
