package main

import (
	"fmt"

	"github.com/rivo/tview"

	"github.com/marcelocantos/flac/internal/data"
	"github.com/marcelocantos/flac/internal/pinyin"
	"github.com/marcelocantos/flac/internal/refdata"
	"github.com/marcelocantos/flac/internal/ui"
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

	db.Populate(rd.WordList().Words)

	root := ui.New(db)
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
			if outcome := Assess(pincache, entries, answer); outcome.IsGood() {
				fmt.Fprintf(root.Results, "[green::b]YES![-::-] %s = %s\n",
					word, outcome.pinyins.ColorString())
				if err := root.Results.Good(word, false); err != nil {
					panic(err)
				}
				if err := setup(); err != nil {
					panic(err)
				}
			} else {
				fmt.Fprintf(root.Results, "[red::b]NO! %s\n", outcome.Correction())
				root.Results.Bad(word, false, &attempt)
			}
		})

	app := tview.NewApplication().SetRoot(root, true).EnableMouse(true)
	if err := app.Run(); err != nil {
		panic(err)
	}

	return nil
}

func main() {
	if err := main2(); err != nil {
		panic(err)
	}
}
