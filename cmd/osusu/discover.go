package main

import (
	"embed"
	"strings"

	"github.com/kkoreilly/osusu/osusu"
	"goki.dev/colors"
	"goki.dev/gi/v2/gi"
	"goki.dev/gi/v2/giv"
	"goki.dev/girl/styles"
	"goki.dev/goosi/events"
	"goki.dev/grows/jsons"
	"goki.dev/grr"
)

//go:embed recipes.json
var recipesFS embed.FS

var recipes []*osusu.Recipe

func configDiscover(rf *gi.Frame, mf *gi.Frame) {
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
		for _, recipe := range recipes {
			grr.Log(recipe.Init())
		}
	}

	for i, recipe := range recipes {
		recipe := recipe

		if i > 20 {
			break
		}

		grr.Log(recipe.CategoryFlag.SetString(strings.Join(recipe.Category, "|")))
		grr.Log(recipe.CuisineFlag.SetString(strings.Join(recipe.Cuisine, "|")))

		if !bitFlagsOverlap(recipe.CategoryFlag, curOptions.Categories) ||
			!bitFlagsOverlap(recipe.CuisineFlag, curOptions.Cuisines) {
			continue
		}

		rc := gi.NewFrame(rf)
		cardStyles(rc)

		img := getImageFromURL(recipe.Image)
		if img != nil {
			gi.NewImage(rc).SetImage(img)
		}

		gi.NewLabel(rc).SetType(gi.LabelHeadlineSmall).SetText(recipe.Name)

		castr := friendlyBitFlagString(recipe.CategoryFlag)
		custr := friendlyBitFlagString(recipe.CuisineFlag)
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

		rc.OnClick(func(e events.Event) {
			addRecipe(rf, recipe, rc, mf)
		})
	}

	rf.Update()
	rf.UpdateEndLayout(updt)
}

func addRecipe(rf *gi.Frame, recipe *osusu.Recipe, rc *gi.Frame, mf *gi.Frame) {
	d := gi.NewBody().AddTitle("Add recipe")
	giv.NewStructView(d).SetStruct(recipe).SetReadOnly(true)
	d.AddBottomBar(func(pw gi.Widget) {
		d.AddCancel(pw)
		d.AddOk(pw).SetText("Add").OnClick(func(e events.Event) {
			meal := &osusu.Meal{
				Name:        recipe.Name,
				Description: recipe.Description,
				Image:       recipe.Image,
				Category:    recipe.CategoryFlag,
				Cuisine:     recipe.CuisineFlag,
			}
			meal.Source.SetFlag(true, osusu.Cooking)
			newMeal(rf, mf, meal)
		})
	})
	d.NewFullDialog(rc).Run()
}
