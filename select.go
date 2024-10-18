package main

import (
	"log"

	fuzzy "github.com/ktr0731/go-fuzzyfinder"
)

type model struct {
	clouds []Clouds
	cursor int
}

func rcselector(c []Clouds) string {
	idx, err := fuzzy.Find(
		c,
		func(i int) string {
			return c[i].name
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	return c[idx].name
}
