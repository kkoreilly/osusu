package main

import (
	"embed"

	"github.com/kkoreilly/osusu/osusu"
	"goki.dev/gi/v2/gi"
	"goki.dev/girl/styles"
	"goki.dev/grows/jsons"
)

//go:embed recipes.json
var recipesFS embed.FS

func configDiscover(rf *gi.Frame) {
	rf.Style(func(s *styles.Style) {
		s.Wrap = true
	})

	recipes := []osusu.Recipe{}
	err := jsons.OpenFS(&recipes, recipesFS, "recipes.json")
	if err != nil {
		gi.ErrorDialog(rf, err)
		return
	}
	for _, recipe := range recipes {
		rc := gi.NewFrame(rf)
		cardStyles(rc)
		gi.NewLabel(rc).SetText(recipe.Name)
	}
}
