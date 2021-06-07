package main

import (
	"fmt"
	"os"

	"github.com/go-errors/errors"
	"github.com/rivo/tview"

	"github.com/marcelocantos/flac/internal/pkg/assess"
	"github.com/marcelocantos/flac/internal/pkg/data"
	"github.com/marcelocantos/flac/internal/pkg/pinyin"
	"github.com/marcelocantos/flac/internal/pkg/refdata"
	"github.com/marcelocantos/flac/internal/pkg/ui"
)

func main2() error {
	pincache := pinyin.Cache{}
	rd, err := refdata.New()
	if err != nil {
		return err
	}
	_ = rd

	db, err := data.NewDatabase("flac.db")
	if err != nil {
		return err
	}

	db.Populate(rd.WordList.Words)

	root := ui.New(db, rd)
	var word string
	var attempt int

	setup := func() error {
		var err error
		word, err = db.HeadWord()
		if err != nil {
			return err
		}
		root.Input.SetLabel(word + ":")
		root.Input.SetText("")
		attempt = 1
		return nil
	}

	if err := setup(); err != nil {
		return err
	}

	root.Input.
		SetValidSyllables(rd.Dict.Syllables).
		SetSubmit(func(answer string) {
			root.Input.SetText("")
			entries := rd.Dict.Entries[word]
			if entries == nil {
				panic("no entry for " + word)
			}
			if outcome := assess.Assess(pincache, word, entries, answer); outcome.IsGood() {
				if err := root.Results.Good(word, false); err != nil {
					panic(err)
				}
				if err := setup(); err != nil {
					panic(err)
				}
			} else {
				root.Results.Bad(word, outcome, false, &attempt)
			}
		})

	app := tview.NewApplication().SetRoot(root, true).EnableMouse(true)
	if err := app.Run(); err != nil {
		panic(err)
	}

	return nil
}

func main() {
	defer func() {
		if err, is := recover().(*errors.Error); is {
			fmt.Fprintln(os.Stderr, err.ErrorStack())
			fmt.Println(err)
			os.Exit(2)
		}
	}()
	if err := main2(); err != nil {
		if err, is := err.(*errors.Error); is {
			fmt.Fprintln(os.Stderr, err.ErrorStack())
		}
		fmt.Println(err)
		os.Exit(2)
	}
}
