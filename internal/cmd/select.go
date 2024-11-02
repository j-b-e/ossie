package cmd

import (
	"log"
	"strings"

	"github.com/j-b-e/ossie/internal/model"
	fuzzy "github.com/ktr0731/go-fuzzyfinder"
)

func preview(c model.Cloud) string {
	var builder strings.Builder
	builder.WriteString("Name: " + c.Name + "\n")
	for k, v := range c.Env {
		if k == "OS_PASSWORD" {
			v = "****"
		}
		builder.WriteString(k + ": " + v + "\n")
	}
	return builder.String() + "\n"
}

func selector(c model.Clouds) model.Cloud {
	idx, err := fuzzy.Find(
		c,
		func(i int) string {
			return c[i].Name
		},
		fuzzy.WithPreviewWindow(func(i int, _ int, _ int) string {
			if i == -1 {
				return "Ossie cantz reed."
			}
			return preview(c[i])
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	return c[idx]
}
