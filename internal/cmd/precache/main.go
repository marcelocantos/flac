package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/go-errors/errors"
	"github.com/spf13/afero"
)

func main2() error {
	if len(os.Args) != 5 || os.Args[1] != "-o" {
		return errors.Errorf("usage: %s -o dest words dict:dict:... ", path.Base(os.Args[0]))
	}
	dest := os.Args[2]
	wordsPath := os.Args[3]
	dictPaths := strings.Split(os.Args[4], ":")

	return cacheRefData(afero.NewOsFs(), wordsPath, dictPaths, dest)
}

func main() {
	defer func() {
		r := recover()
		switch err := r.(type) {
		case nil:
			return
		case *errors.Error:
			fmt.Fprintln(os.Stderr, err.ErrorStack())
		case error:
		default:
			r = errors.Wrap(err, 0)
		}
		fmt.Printf("panic: %s\n", r)
		os.Exit(2)
	}()
	if err := main2(); err != nil {
		if err, is := err.(*errors.Error); is {
			fmt.Fprintln(os.Stderr, err.ErrorStack())
		}
		fmt.Println(err)
		os.Exit(2)
	}
}
