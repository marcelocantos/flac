package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rivo/tview"

	"github.com/marcelocantos/flac/internal/pkg/assess"
	"github.com/marcelocantos/flac/internal/pkg/data"
	"github.com/marcelocantos/flac/internal/pkg/mainwrap"
	"github.com/marcelocantos/flac/internal/pkg/outcome"
	"github.com/marcelocantos/flac/internal/pkg/proto/refdata_pb"
	"github.com/marcelocantos/flac/internal/pkg/refdata"
	"github.com/marcelocantos/flac/internal/pkg/ui"
)

func main2() (err error) {
	phrase := flag.String("phrase", "", "Focus this session on words from a phrase")
	words := flag.String("words", "", "Focus this session on words from a comma-separated list")
	flag.Parse()

	if *phrase != "" && *words != "" {
		fmt.Println("-phrase and -words are mutually exclusive")
		os.Exit(1)
	}

	if err = hackTerminfo(); err != nil {
		return err
	}

	rd, err := refdata.New()
	if err != nil {
		return err
	}

	db, err := data.NewDatabase("flac.db")
	if err != nil {
		return err
	}

	var focus string
	var focusWords []string
	if *phrase != "" {
		focus = "phrase:" + *phrase
		focusWords, err = parsePhrase(*phrase, rd)
		if err != nil {
			return err
		}
	} else if *words != "" {
		focus = "words:" + *words
		focusWords = parseWords(*words)
	} else {
		focusWords = rd.WordList.Words
	}

	if err := db.Populate(focus, focusWords); err != nil {
		return err
	}

	root := ui.New(db, rd)
	var word string
	var variantDef string
	var entries *refdata_pb.CEDict_Entries
	var attempt int

	setup := func() error {
		var err error
		word, err = db.HeadWord()
		if err != nil {
			return err
		}
		score, err := db.WordScore(word)
		switch err := err.(type) {
		case data.ErrNotFound:
			score = 0
		case nil:
		default:
			return err
		}
		labelWidth := root.Answer.SetWord(word, score)

		variantDef, entries = refdata.RandomDefinition(word, rd.Dict.Entries[word])

		if variantDef == "" {
			root.ResizeItem(root.Hint, 0, 0)
		} else {
			root.Hint.SetText(strings.Repeat(" ", labelWidth) + ui.DecorateDefinition(variantDef))
			root.ResizeItem(root.Hint, 1, 0)
		}

		root.Answer.SetText("")

		attempt = 1
		return nil
	}

	if err := setup(); err != nil {
		return err
	}

	root.Answer.
		SetExitFunc(func() {
			panic(mainwrap.Stop)
		}).
		SetGiveUpFunc(func() {
			outcome := &outcome.Outcome{
				Word:    word,
				Entries: entries,
			}
			if outcome.Entries == nil {
				panic("no entry for " + word)
			}
			if err := root.Results.GiveUp(outcome); err != nil {
				panic(err)
			}
		}).
		SetSubmitFunc(func(answer string) {
			root.Answer.SetText("")
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
		})
	app := tview.NewApplication().SetRoot(root, true)
	root.Answer.App = app
	if err := app.Run(); err != nil {
		return err
	}

	return nil
}

func main() {
	mainwrap.Main(main2)
}
