package main

import (
	"github.com/rivo/tview"
	"github.com/spf13/afero"

	"github.com/marcelocantos/flac/internal/data"
	"github.com/marcelocantos/flac/internal/refdata"
	"github.com/marcelocantos/flac/internal/ui"
)

func main2() error {
	refdata, err := refdata.New(
		afero.NewBasePathFs(afero.NewOsFs(), "refdata"))
	if err != nil {
		return err
	}
	_ = refdata

	db, err := data.NewDatabase("flac.db")
	if err != nil {
		return err
	}

	db.Populate(refdata.Words())

	root := ui.New()
	headWord, err := db.HeadWord()
	if err != nil {
		return err
	}
	root.Input.SetLabel(headWord + ":")

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
