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

	root := ui.New()
	headWord, err := db.HeadWord()
	if err != nil {
		return err
	}
	root.Input.
		SetValidSyllables(rd.Dict.Syllables).
		SetSubmit(func(answer string) {
			root.Input.SetText("")
			entries := rd.Dict.Entries[headWord].GetDefinitions()
			if entries != nil && entries[answer] != nil {
				fmt.Fprintf(root.Results, "[green::b]YES![-::-] %s = %s\n",
					headWord, pincache.MustPinyin(answer).ColorString())
			} else {
				fmt.Fprintf(root.Results, "[red::b]NO!\n")
			}
		}).
		SetLabel(headWord + ":")

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
