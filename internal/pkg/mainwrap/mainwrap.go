package mainwrap

import (
	"fmt"
	"os"

	"github.com/go-errors/errors"
)

type StopT struct{}

var Stop = StopT{}

func Main(main func() error) {
	defer func() {
		if r := recover(); r != nil {
			if _, is := r.(StopT); !is {
				panic(r)
			}
		}
	}()
	if err := main(); err != nil {
		if err, is := err.(*errors.Error); is {
			fmt.Fprintln(os.Stderr, err.ErrorStack())
		}
		fmt.Println(err)
		os.Exit(2)
	}
}
