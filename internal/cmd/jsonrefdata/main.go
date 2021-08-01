package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/marcelocantos/flac/internal/pkg/mainwrap"
	"github.com/marcelocantos/flac/internal/pkg/refdata"
	"google.golang.org/protobuf/encoding/protojson"
)

func main2() error {
	from := os.Args[1]
	to := os.Args[2]

	log.Printf("Converting %s to %s...", from, to)

	data, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}

	rd, err := refdata.NewFromBytes(data)
	if err != nil {
		return err
	}

	marshaler := protojson.MarshalOptions{
		// Multiline: true,
		// Indent:    "  ",
	}
	json, err := marshaler.Marshal(rd)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(to, json, 0777); err != nil {
		return err
	}

	log.Print("Done")
	return nil
}

func main() {
	mainwrap.Main(main2)
}
