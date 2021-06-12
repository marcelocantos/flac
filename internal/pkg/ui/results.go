package ui

import (
	"fmt"
	"math"
	"math/rand"
	"strings"

	"github.com/rivo/tview"

	"github.com/marcelocantos/flac/internal/pkg/assess"
	"github.com/marcelocantos/flac/internal/pkg/data"
	"github.com/marcelocantos/flac/internal/pkg/proto/refdata"
)

var brailleBars = []string{"", "⡀", "⡄", "⡆", "⡇", "⣇", "⣧", "⣷"}

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

	db         *data.Database
	rd         *refdata.RefData
	wordScores map[string]int
	history    []string
	goods      []string

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
		wordScores:   map[string]int{},
		scoreChanged: func(word string, score int) {},
	}
	return results.refreshText()
}

func (r *Results) refreshText() *Results {
	r.SetText("")
	fmt.Fprintf(r, "%s你好，一起学中文吧！\n", strings.Repeat("\n", 999))

	// Abuse history as a preallocated buffer for output.
	output := append(r.history, r.goodsReport()...)
	r.history = output[:len(r.history)]

	for i, h := range output {
		if i == len(r.history)-1 {
			h = strings.SplitN(h, "\034", 2)[0]
		}
		fmt.Fprintf(r, "\n%s", h)
	}
	return r
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
	r.trimEphemeralContent()
	r.history = append(r.history, lines...)
	r.refreshText()
}

func (r *Results) trimEphemeralContent(line ...string) {
	if len(r.history) > 0 {
		last := len(r.history) - 1
		h := strings.SplitN(r.history[last], "\034", 2)
		r.history[last] = h[len(h)-1]
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

func (r *Results) SetScoreChanged(f func(word string, score int)) *Results {
	r.scoreChanged = f
	return r
}

func (r *Results) Good(word string, easy bool) error {
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

	r.goods = append(r.goods, word+brailleScore(score))
	r.refreshText()

	return nil
}

func (r *Results) Bad(word string, outcome *assess.Outcome, easy bool, attempt *int) error {
	penalty := math.Sqrt(float64(1 + *attempt))
	*attempt++

	if err := r.bump(word, func(score int) (int, bool) {
		// Multiply score by 1/2√(1 + attempt).
		return atLeast(1)(score / int(2*penalty)), false
	}); err != nil {
		return err
	}

	r.appendHistory(r.goodsReport()...)
	r.goods = nil

	r.appendHistory(outcome.Correction())

	return nil
}

func (r *Results) Skip(word string, easy bool, attempt int) error {
	r.trimEphemeralContent()

	return r.bump(word, func(score int) (int, bool) {
		return atLeast(1)(score / 8), true
	})
}
