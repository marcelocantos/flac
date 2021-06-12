package ui

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var expectedInputRE = regexp.MustCompile(`^(.*)\((.*)\)(.*)$`)

func TestInputREEmpty(t *testing.T) {
	t.Parallel()

	assertMatch(t, "()")
}

func TestInputREBadStart(t *testing.T) {
	t.Parallel()

	assertNoMatch(t, "1")
	assertNoMatch(t, "12")
	assertNoMatch(t, "/")
	assertNoMatch(t, "!")
}

func TestInputREGoodStart(t *testing.T) {
	t.Parallel()

	assertMatch(t, "(a)")
	assertMatch(t, "(a1)")
	assertMatch(t, "(a2) ")
	assertMatch(t, "(a3)      ")
	assertMatch(t, "(a4)")
	assertMatch(t, "(a5)")
	assertMatch(t, "(a12)")
	assertMatch(t, "(yi)")
	assertMatch(t, "(yi1)")
	assertMatch(t, "(yi12)")
	assertMatch(t, "(zhuang12345)")
}

func TestInputREBadTones(t *testing.T) {
	t.Parallel()

	assertNoMatch(t, "a0")
	assertNoMatch(t, "a6")
	assertNoMatch(t, "a9")
	assertNoMatch(t, "yi0")
	assertNoMatch(t, "yi6")
	assertNoMatch(t, "yi9")
	assertNoMatch(t, "zhuang0123456789")
}

func TestInputREBadFirstChar(t *testing.T) {
	t.Parallel()

	assertNoMatch(t, "yi/")
	assertNoMatch(t, "12/")
	assertNoMatch(t, "yi9/")
}

func TestInputREGoodFirstChar(t *testing.T) {
	t.Parallel()

	assertMatch(t, "yi1/()")
	assertMatch(t, "yi12 / ()")
	assertMatch(t, "yi1 / (s)   ")
	assertMatch(t, "yi1/(shi)")
	assertMatch(t, "yi1/(shi4)")
}

func TestInputRESeveralChars(t *testing.T) {
	t.Parallel()

	assertMatch(t, "jiang1/()")
	assertMatch(t, "jiang14/()")
	assertMatch(t, "jiang14/(q)")
	assertMatch(t, "jiang14/(qiang)")
	assertMatch(t, "jiang14/(qiang1)")
	assertMatch(t, "yi1ding1(bu4)")
	assertMatch(t, "yi1ding1bu4(s)")
	assertMatch(t, "yi1ding1bu4 (s)")
	assertMatch(t, "yi1 ding1 bu4 (s)")
	assertMatch(t, "yi1 ding1 bu4 (shi)")
	assertMatch(t, "yi1 ding1 bu4 (shi2)")

	assertMatch(t, "nuo2(na)")
	assertMatch(t, "nuo2(na1)")
	assertMatch(t, "nuo2(na13)")
	assertMatch(t, "nuo2(na134)")

	assertMatch(t, "Ya4 dang1 ·()")
	assertMatch(t, "Ya4 dang1 · ()")
	assertMatch(t, "Ya4 dang1 · (S)")
	assertMatch(t, "Ya4 dang1 · (Si1)")
	assertMatch(t, "Ya4 dang1 · (Si1) ")
	assertMatch(t, "Ya4 dang1 · Si1 (mi4)")
	assertMatch(t, "Ya4 dang1 · Si1 (mi4)")

	assertMatch(t, "yi1 bu4 zuo4,()")
	assertMatch(t, "yi1 bu4 zuo4, ()")
	assertMatch(t, "yi1 bu4 zuo4 , ()")
	assertMatch(t, "yi1 bu4 zuo4 , (e)")
	assertMatch(t, "yi1 bu4 zuo4 , (er4)")
	assertMatch(t, "yi1 bu4 zuo4 , (er4) ")
	assertMatch(t, "yi1 bu4 zuo4 , er4 bu4 (xiu1)")
}

func TestInputREBadSeveralChars(t *testing.T) {
	t.Parallel()

	assertNoMatch(t, "1 ding1 bu4 shi2")
	assertNoMatch(t, "yi1 ding bu4 shi2")
	assertNoMatch(t, "yi1 ding1 bu4 shi0")

	assertNoMatch(t, "yi1 bu4 zuo,")
}

// assertMatch takes patterns of the form "prefix(lastchar)", asserting that
// inputRE matches prefix+lastchar and captures lastchar as group 1.
func assertMatch(t *testing.T, pattern string) bool {
	t.Helper()

	e := expectedInputRE.FindStringSubmatch(pattern)
	require.NotEmpty(t, e)
	all := e[1] + e[2] + e[3]
	expected := e[2]

	m := inputRE.FindStringSubmatch(all)
	return assert.NotNil(t, m) && assert.Equal(t, expected, m[1])
}

func assertNoMatch(t *testing.T, input string) bool {
	t.Helper()

	m := inputRE.FindStringSubmatch(input)
	return assert.Nil(t, m)
}
