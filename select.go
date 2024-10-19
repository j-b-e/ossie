package main

import (
	"log"
	"strings"

	fuzzy "github.com/ktr0731/go-fuzzyfinder"
)

func envStr(c Cloud, env, name string) string {
	if env, ok := c.Env[env]; ok {
		return name + ": " + env
	}
	return ""
}

func rcselector(c []Cloud) Cloud {
	idx, err := fuzzy.Find(
		c,
		func(i int) string {
			return c[i].Name
		},
		fuzzy.WithPreviewWindow(func(i int, width int, height int) string {
			if i == -1 {
				return "Ossie cantz reed."
			}
			return strings.Join([]string{
				"Name: " + c[i].Name,
				envStr(c[i], "OS_USERNAME", "Username"),
				envStr(c[i], "OS_DOMAIN_NAME", "Domain"),
				envStr(c[i], "OS_REGION_NAME", "Region"),
			}, "\n")
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	return c[idx]
}
