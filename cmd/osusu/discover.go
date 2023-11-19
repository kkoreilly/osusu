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

var recipes []osusu.Recipe

func configDiscover(rf *gi.Frame) {
	if rf.HasChildren() {
		rf.DeleteChildren(true)
	}
	updt := rf.UpdateStart()

	rf.Style(func(s *styles.Style) {
		s.Wrap = true
	})

	if recipes == nil {
		err := jsons.OpenFS(&recipes, recipesFS, "recipes.json")
		if err != nil {
			gi.ErrorDialog(rf, err)
			return
		}
	}

	for i, recipe := range recipes {
		if i > 20 {
			break
		}

		var categories osusu.Categories
		grr.Log0(categories.SetString(strings.Join(recipe.Category, "|")))

		var cuisines osusu.Cuisines
		grr.Log0(cuisines.SetString(strings.Join(recipe.Cuisine, "|")))

		if !bitFlagsOverlap(categories, curOptions.Categories) ||
			!bitFlagsOverlap(cuisines, curOptions.Cuisines) {
			continue
		}

		rc := gi.NewFrame(rf)
		cardStyles(rc)

		img := getImageFromURL(recipe.Image)
		if img != nil {
			gi.NewImage(rc).SetImage(img, 0, 0)
		}

		gi.NewLabel(rc).SetType(gi.LabelHeadlineSmall).SetText(recipe.Name)

		castr := friendlyBitFlagString(categories)
		custr := friendlyBitFlagString(cuisines)
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

	rf.Update()
	rf.UpdateEndLayout(updt)
}
