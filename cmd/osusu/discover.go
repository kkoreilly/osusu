package main

import (
	"embed"
	"strings"

	"github.com/kkoreilly/osusu/osusu"
	"goki.dev/colors"
	"goki.dev/gi/v2/gi"
	"goki.dev/girl/styles"
	"goki.dev/grows/jsons"
	"goki.dev/grr"
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
	for i, recipe := range recipes {
		if i > 20 {
			break
		}

		rc := gi.NewFrame(rf)
		cardStyles(rc)

		img := getImageFromURL(recipe.Image)
		if img != nil {
			gi.NewImage(rc).SetImage(img, 0, 0)
		}

		gi.NewLabel(rc).SetType(gi.LabelHeadlineSmall).SetText(recipe.Name)

		var ca osusu.Categories
		grr.Log0(ca.SetString(strings.Join(recipe.Category, "|")))
		castr := friendlyBitFlagString(ca)

		var cu osusu.Cuisines
		grr.Log0(cu.SetString(strings.Join(recipe.Cuisine, "|")))
		custr := friendlyBitFlagString(cu)
		text := castr
		if castr != "" && custr != "" {
			text += " • "
		}
		text += custr
		gi.NewLabel(rc).SetText(text).Style(func(s *styles.Style) {
			s.Color = colors.Scheme.OnSurfaceVariant
		})

		recipe.Score.ComputeTotal(curOptions)
		scoreGrid(rc, &recipe.Score, true)
	}
}
