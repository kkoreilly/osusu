package main

import (
	"context"
	"embed"
	"fmt"
	"strings"

	"github.com/kkoreilly/osusu/osusu"
	"github.com/kkoreilly/osusu/otextencoding"
	"github.com/nlpodyssey/cybertron/pkg/models/bert"
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

//go:embed textEncodingVectors.json
var textEncodingVectorsFS embed.FS

var textEncodingVectors map[string][]float32

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
			gi.ErrorDialog(rf, err, "Error opening recipes")
			return
		}
		for _, recipe := range recipes {
			grr.Log(recipe.Init())
		}

		err = jsons.OpenFS(&textEncodingVectors, textEncodingVectorsFS, "textEncodingVectors.json")
		if err != nil {
			gi.ErrorDialog(rf, err, "Error opening recipe text encoding vectors")
			return
		}

		err = otextencoding.LoadModel()
		if err != nil {
			gi.ErrorDialog(rf, err, "Error loading text encoding model")
			return
		}
	}

	var meals []*osusu.Meal
	err := osusu.DB.Find(&meals).Error
	if err != nil {
		gi.ErrorDialog(rf, err).Run()
	}
	for _, meal := range meals {
		meal := meal

		res, err := otextencoding.Model.Encode(context.TODO(), meal.Text(), int(bert.MeanPooling))
		if err != nil {
			gi.ErrorDialog(rf, err, "Error text encoding meal")
			continue
		}
		fmt.Println(meal.Name, res.Vector.Data().Len())
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

		if recipe.Image != "" {
			img := gi.NewImage(rc)
			go func() {
				img.SetImage(getImageFromURL(recipe.Image))
				img.Update()
			}()
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
