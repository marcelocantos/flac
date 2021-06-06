package ui

import (
	"fmt"
	"math"
	"math/rand"
	"strings"

	"github.com/rivo/tview"

	"github.com/marcelocantos/flac/internal/pkg/assess"
	"github.com/marcelocantos/flac/internal/pkg/data"
	"github.com/marcelocantos/flac/internal/pkg/pinyin"
	"github.com/marcelocantos/flac/internal/pkg/proto/refdata"
)

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

	db        *data.Database
	rd        *refdata.RefData
	pincache  pinyin.Cache
	wordsSeen map[string]bool
	history   []string
	goods     []string

	// Handlers
	scoreChanged func(word string, score int)
}

func newResults(db *data.Database, rd *refdata.RefData) *Results {
	view := tview.NewTextView()
	view.SetDynamicColors(true)
	view.SetBorder(true)
	view.SetTitle("flac: learn 中文")

	results := &Results{
		TextView:     view,
		db:           db,
		rd:           rd,
		pincache:     pinyin.Cache{},
		wordsSeen:    map[string]bool{},
		scoreChanged: func(word string, score int) {},
	}
	return results.refreshText()
}

func (r *Results) refreshText() *Results {
	r.SetText("")
	fmt.Fprintf(r, "%s你好！", strings.Repeat("\n", 999))
	for _, h := range append(r.history, r.goodsReport()...) {
		fmt.Fprintf(r, "\n%s", h)
	}
	return r
}

func (r *Results) goodsReport() []string {
	if r.goods == nil {
		return nil
	}
	return []string{
		fmt.Sprintf("[green::b]%s[-::-]", strings.Join(r.goods, " ")),
	}
}

func (r *Results) appendHistory(lines ...string) {
	r.history = append(r.history, lines...)
	r.refreshText()
}

func (r *Results) bump(word string, bump func(score int) (int, bool)) error {
	r.wordsSeen[word] = true
	score, err := r.score(word)
	if err != nil {
		return err
	}
	score, newpos := bump(score)
	pos := -1
	if newpos {
		pos = score + rand.Intn(1+score*3/2-score)
	}
	return r.db.UpdateScoreAndPos(word, score, pos)
}

func (r *Results) score(word string) (int, error) {
	score, err := r.db.WordScore(word)
	if _, is := err.(data.ErrNotFound); is {
		return -1, nil
	}
	return score, err
}

func (r *Results) SetScoreChanged(f func(word string, score int)) *Results {
	r.scoreChanged = f
	return r
}

func (r *Results) Good(word string, easy bool) error {
	if err := r.bump(word, func(score int) (int, bool) {
		return atLeast(2)(score) * 3 / 2, true
	}); err != nil {
		return err
	}

	r.goods = append(r.goods, word)
	r.refreshText()
	return nil
}

func (r *Results) Bad(word string, outcome *assess.Outcome, easy bool, attempt *int) error {
	// Multiply score by 2/(3*sqrt(attempt)).
	penalty := int(10 * math.Sqrt(float64(*attempt)))
	*attempt++

	if err := r.bump(word, func(score int) (int, bool) {
		return atLeast(1)(score * 20 / 3 / penalty), false
	}); err != nil {
		return err
	}

	r.appendHistory(r.goodsReport()...)
	r.goods = nil

	r.appendHistory(fmt.Sprintf("❌ %s", outcome.Correction()))

	return nil
}

func (r *Results) Skip(word string, easy bool, attempt int) error {
	return r.bump(word, func(score int) (int, bool) {
		return atLeast(1)(score / 8), true
	})
}
