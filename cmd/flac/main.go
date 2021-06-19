package main

import (
	"fmt"
	"os"

	"github.com/go-errors/errors"
	"github.com/rivo/tview"

	"github.com/marcelocantos/flac/internal/pkg/assess"
	"github.com/marcelocantos/flac/internal/pkg/data"
	"github.com/marcelocantos/flac/internal/pkg/outcome"
	"github.com/marcelocantos/flac/internal/pkg/refdata"
	"github.com/marcelocantos/flac/internal/pkg/ui"
)

func main2() (err error) {
	if err = hackTerminfo(); err != nil {
		return err
	}

	rd, err := refdata.New()
	if err != nil {
		return err
	}
	_ = rd

	db, err := data.NewDatabase("flac.db")
	if err != nil {
		return err
	}

	if err := db.Populate(rd.WordList); err != nil {
		return err
	}

	root := ui.New(db, rd)
	var word string
	var attempt int

	setup := func() error {
		var err error
		word, err = db.HeadWord()
		if err != nil {
			return err
		}
		root.Input.SetWord(word)
		root.Input.SetText("")
		attempt = 1
		return nil
	}

	if err := setup(); err != nil {
		return err
	}

	root.Input.
		SetValidSyllables(rd.Dict.ValidSyllables).
		SetSubmitFunc(func(answer string) {
			root.Input.SetText("")
			entries := rd.Dict.Entries[word]
			if entries == nil {
				panic(errors.Errorf("no entry for %s", word))
			}
			outcome := assess.Assess(word, entries, answer)
			if outcome.Pass() {
				if err := root.Results.Good(word, outcome, false); err != nil {
					panic(err)
				}
				if err := setup(); err != nil {
					panic(err)
				}
			} else {
				if err := root.Results.NotGood(outcome, false, &attempt); err != nil {
					panic(err)
				}
			}
		}).
		SetGiveUpFunc(func() {
			outcome := &outcome.Outcome{
				Word:    word,
				Entries: rd.Dict.Entries[word],
			}
			if outcome.Entries == nil {
				panic("no entry for " + word)
			}
			if err := root.Results.GiveUp(outcome); err != nil {
				panic(err)
			}
		}).
		SetChangedFunc(func(text string) {
			if text != "" {
				root.Results.BlankOutMessages()
			}
		})
	app := tview.NewApplication().SetRoot(root, true)
	root.Input.App = app
	if err := app.Run(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := main2(); err != nil {
		if err, is := err.(*errors.Error); is {
			fmt.Fprintln(os.Stderr, err.ErrorStack())
		}
		fmt.Println(err)
		os.Exit(2)
	}
}
