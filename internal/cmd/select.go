package cmd

import (
	"github.com/j-b-e/ossie/internal/model"
	fuzzy "github.com/ktr0731/go-fuzzyfinder"
)

func selector(c model.Clouds) (model.Cloud, error) {
	idx, err := fuzzy.Find(
		c,
		func(i int) string {
			return c[i].Name
		},
		fuzzy.WithPreviewWindow(func(i int, _ int, _ int) string {
			if i == -1 {
				return "Ossie cantz reed."
			}
			return c[i].String()
		}),
	)
	if err != nil {
		return model.Cloud{}, err
	}
	return c[idx], nil
}
