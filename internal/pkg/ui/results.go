package ui

import (
	"math"
	"math/rand"

	"github.com/rivo/tview"

	"github.com/marcelocantos/flac/internal/pkg/data"
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
	wordsSeen map[string]bool

	// Handlers
	scoreChanged func(word string, score int)
}

func newResults(db *data.Database) *Results {
	view := tview.NewTextView()
	view.SetDynamicColors(true)
	view.SetBorder(true)
	view.SetTitle("flac: learn 中文")

	return &Results{
		TextView:     view,
		db:           db,
		wordsSeen:    map[string]bool{},
		scoreChanged: func(word string, score int) {},
	}
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
	return r.bump(word, func(score int) (int, bool) {
		return atLeast(2)(score) * 3 / 2, true
	})
}

func (r *Results) Bad(word string, easy bool, attempt *int) error {
	// Multiply score by 2/(3*sqrt(attempt)).
	penalty := int(10 * math.Sqrt(float64(*attempt)))
	*attempt++

	return r.bump(word, func(score int) (int, bool) {
		return atLeast(1)(score * 20 / 3 / penalty), false
	})
}

func (r *Results) Skip(word string, easy bool, attempt int) error {
	return r.bump(word, func(score int) (int, bool) {
		return atLeast(1)(score / 8), true
	})
}
