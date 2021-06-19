package ui

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"strings"

	"github.com/rivo/tview"

	"github.com/marcelocantos/flac/internal/pkg/data"
	"github.com/marcelocantos/flac/internal/pkg/outcome"
	"github.com/marcelocantos/flac/internal/pkg/pinyin"
	"github.com/marcelocantos/flac/internal/pkg/proto/refdata"
)

var (
	brailleBars = []string{"", "⡀", "⡄", "⡆", "⡇", "⣇", "⣧", "⣷"}
	superDigits = []rune{'⁰', '¹', '²', '³', '⁴', '⁵', '⁶', '⁷', '⁸', '⁹'}
)

func logscore(score int) float64 {
	return math.Log(float64(score)) / math.Log(2)
}

func brailleScore(score int) string {
	if score <= 0 {
		return ""
	}
	s := int(logscore(score))
	return strings.Repeat("⣿", s/8) + brailleBars[s%8]
}

func superNumber(n int) string {
	s := strconv.Itoa(n)
	var sb strings.Builder
	for _, d := range s {
		sb.WriteRune(superDigits[d-'0'])
	}
	return sb.String()
}

func atLeast(min int) func(i int) int {
	return func(i int) int {
		if i < min {
			return min
		}
		return i
	}
}

type Results struct {
	*tview.TextView

	db *data.Database
	rd *refdata.RefData

	wordScores map[string]int

	stale        bool
	refreshCount int

	history []string
	goods   []string
	msgs    messages

	scoreChangedFunc func(word string, score int)
}

func newResults(db *data.Database, rd *refdata.RefData) *Results {
	view := tview.NewTextView()
	view.SetDynamicColors(true)
	view.SetBorder(true)
	view.SetTitle("flac: learn 中文")

	r := &Results{
		TextView:         view,
		db:               db,
		rd:               rd,
		wordScores:       map[string]int{},
		scoreChangedFunc: func(word string, score int) {},
		stale:            true,
	}

	r.refresh()()

	return r
}

func (r *Results) SetScoreChangedFunc(f func(word string, score int)) *Results {
	r.scoreChangedFunc = f
	return r
}

func (r *Results) Good(word string, o *outcome.Outcome, easy bool) error {
	defer r.refresh()()

	if err := r.bump(word, func(score int) (int, bool) {
		return atLeast(2)(2 * score), true
	}); err != nil {
		return err
	}

	score, err := r.score(word)
	if err != nil {
		return err
	}

	r.trimEphemeralContent()
	r.appendGoods(word + brailleScore(score))
	r.ClearMessages()

	maxPrefixLen := 0
	maxDigits := 0
	for word, entry := range o.Entries.Definitions {
		pword := pinyin.MustNewWord(word)
		prefixLen := len([]rune(pword.String())) + 4
		if maxPrefixLen < prefixLen {
			maxPrefixLen = prefixLen
		}
		digits := len([]rune(superNumber(len(entry.Definitions) - 1)))
		if maxDigits < digits {
			maxDigits = digits
		}
	}
	prefix := "\n" + strings.Repeat(" ", maxPrefixLen)

	for word, entry := range o.Entries.Definitions {
		var sb strings.Builder
		pword := pinyin.MustNewWord(word)
		prefixLen := len([]rune(pword.String())) + 4
		fmt.Fprintf(&sb, "%s%s 👉 ", pword.ColorString(), strings.Repeat(" ", maxPrefixLen-prefixLen))
		for i, def := range entry.Definitions {
			if i > 0 {
				sb.WriteString(prefix)
			}
			fmt.Fprintf(&sb, "[#909090::]%*s[-::]%s", maxDigits, superNumber(i+1), accentPhrase(def))
		}
		r.appendMessage("%s", sb.String())
	}

	return nil
}

func (r *Results) NotGood(o *outcome.Outcome, easy bool, attempt *int) error {
	defer r.refresh()()

	if o.Fail() {
		if err := r.bad(o, easy, attempt); err != nil {
			return err
		}
	}

	r.ClearMessages()

	if len(o.Bad) > 0 {
		prefix := strings.Repeat(" ", 3+2*len([]rune(o.Word))+2)
		top := prefix
		var corrections [][]string
		for _, word := range o.Bad {
			// word := o.Bad[i]
			wordLen := len([]rune(word.String()))
			middle := (wordLen - 1) / 2
			tail := wordLen - middle - 1
			var correction string
			if dancis, has := r.rd.Dict.PinyinToSimplified[word.RawString()]; has {
				correction = strings.Join(dancis.Words, " ")
			} else {
				correction = "∅"
			}
			top = fmt.Sprintf("%s %s┬%s", top, strings.Repeat("─", middle), strings.Repeat("─", tail))
			corrections = append(corrections, []string{
				fmt.Sprintf("%s %*s╘👉 %s", prefix, middle, "", correction),
			})
			prefix = fmt.Sprintf("%s %s│%s", prefix, strings.Repeat(" ", middle), strings.Repeat(" ", tail))
		}
		r.appendMessage("[silver::]%s[-::]", top)
		for i := len(corrections) - 1; i >= 0; i-- {
			for _, line := range corrections[i] {
				r.appendMessage("[silver::]%s[-::]", line)
			}
		}

		r.appendHistory(fmt.Sprintf(
			"❌ %s ≠ %s\034❌ [#999999::]%[1]s ≠ [#999999::d]%[3]s[-::-]",
			o.Word, o.Bad.ColorString(), o.Bad.String()))
	}
	if len(o.TooShort) > 0 {
		r.appendMessage("⚠️  Missing characters: %s...", o.TooShort.ColorString())
	}
	if len(o.Bad) == 0 && o.Missing > 0 {
		r.appendMessage("⚠️  Missing alternative%s[-::]", pluralS(o.Missing))
	}
	if len(o.BadTones) > 0 {
		r.appendMessage("[:silver:]🎵[:-:] Only tone(s) need correcting!")
	}

	return nil
}

func (r *Results) GiveUp(outcome *outcome.Outcome) error {
	defer r.refresh()()

	r.trimEphemeralContent()

	r.setMessages(outcome.Correction())
	return r.bump(outcome.Word, func(score int) (int, bool) {
		return atLeast(1)(score / 8), true
	})
}

func (r *Results) taint() {
	r.stale = true
}

func (r *Results) refresh() func() {
	r.refreshCount++
	return func() {
		if r.refreshCount--; r.refreshCount != 0 {
			return
		}
		if r.stale {
			r.SetText("")
			fmt.Fprintf(r, "%s欢迎来到flac，一起学中文吧！\n", strings.Repeat("\n", 999))

			// Abuse history as a preallocated buffer for output.
			output := append(r.history, r.goodsReport()...)
			output = append(output, r.msgs...)
			r.history = output[:len(r.history)]

			for _, h := range output {
				if index := strings.IndexRune(h, '\034'); index >= 0 {
					h = h[:index]
				}
				fmt.Fprintf(r, "\n%s", h)
			}

			r.stale = false
		}
	}
}

func (r *Results) appendGoods(goods ...string) {
	if len(goods) > 0 {
		r.goods = append(r.goods, goods...)
		r.taint()
	}
}

func (r *Results) clearGoods(goods ...string) {
	if len(r.goods) > 0 {
		r.appendHistory(r.goodsReport()...)
		r.goods = nil
		r.taint()
	}
}

func (r *Results) goodsReport() []string {
	if len(r.goods) == 0 {
		return nil
	}
	return []string{
		fmt.Sprintf("[green::b]%s[-::-]", strings.Join(r.goods, " ")),
	}
}

func (r *Results) appendHistory(lines ...string) {
	if len(lines) > 0 {
		r.history = append(r.history, lines...)
		r.taint()
	}
}

func (r *Results) trimEphemeralContent(line ...string) {
	if len(r.history) > 0 {
		last := len(r.history) - 1
		h := strings.SplitN(r.history[last], "\034", 2)
		s := h[len(h)-1]
		if r.history[last] != s {
			r.history[last] = s
			r.taint()
		}
	}
}

func (r *Results) bump(word string, bump func(score int) (int, bool)) error {
	score, err := r.score(word)
	if err != nil {
		return err
	}
	score, newpos := bump(score)
	pos := -1
	if newpos {
		pos = score + rand.Intn(1+score*3/2-score)
	}

	return r.setScoreAndPos(word, score, pos)
}

func (r *Results) score(word string) (int, error) {
	score, has := r.wordScores[word]
	if !has {
		var err error
		score, err = r.db.WordScore(word)
		if _, is := err.(data.ErrNotFound); is {
			return -1, nil
		}
		if err != nil {
			return 0, err
		}
		r.wordScores[word] = score
	}
	return score, nil
}

func (r *Results) setScoreAndPos(word string, score, pos int) error {
	r.wordScores[word] = score
	return r.db.UpdateScoreAndPos(word, score, pos)
}

func (r *Results) bad(outcome *outcome.Outcome, easy bool, attempt *int) error {
	defer r.clearGoods()

	penalty := math.Sqrt(float64(1 + *attempt))
	*attempt++

	if err := r.bump(outcome.Word, func(score int) (int, bool) {
		// Multiply score by 1/2√(1 + attempt).
		return atLeast(1)(score / int(2*penalty)), false
	}); err != nil {
		return err
	}

	r.trimEphemeralContent()
	r.ClearMessages()

	return nil
}

func (r *Results) ClearMessages() {
	defer r.refresh()()

	r.clearMessages()
}

func (r *Results) setMessages(messages ...string) {
	if len(r.msgs) > 0 || len(messages) > 0 {
		if !reflect.DeepEqual(r.msgs, messages) {
			r.msgs = messages
			r.taint()
		}
	}
}

func (r *Results) appendMessage(format string, args ...interface{}) {
	r.msgs = r.msgs.write(format, args...)
	r.taint()
}

func (r *Results) clearMessages() {
	r.setMessages()
}
