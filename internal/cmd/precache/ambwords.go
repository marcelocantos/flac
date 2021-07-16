package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/marcelocantos/flac/internal/pkg/proto/refdata_pb"
)

func AmbiguousWords(dict *refdata_pb.CEDict) map[string]bool {
	wordses := map[string]map[string]bool{}

	digits := regexp.MustCompile(`\d+`)
	for _, entry := range dict.Entries {
		for word := range entry.Entries {
			chars := strings.Split(digits.ReplaceAllString(strings.ToLower(word), ""), " ")
			for i, b := range chars[1:] {
				a := chars[i]
				if a+b == "shenge" {
					log.Print(word)
				}
				wordses[a+b] = map[string]bool{}
			}
		}
	}

	for _, entry := range dict.Entries {
		for word := range entry.Entries {
			chars := strings.Split(digits.ReplaceAllString(strings.ToLower(word), ""), " ")
			for _, a := range chars {
				m, has := wordses[a]
				if !has {
					m = map[string]bool{}
				}
				m[a] = true
			}
			for i, b := range chars[1:] {
				a := chars[i]
				c := a + b
				m, has := wordses[c]
				if !has {
					m = map[string]bool{}
				}
				m[a+" "+b] = true
			}
		}
	}
	ambWords := map[string]bool{}
	for amb, words := range wordses {
		if len(words) > 1 {
			log.Print(amb, " ", words)
			ambWords[amb] = true
		}
	}
	return ambWords
}
