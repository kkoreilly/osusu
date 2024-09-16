package main

import (
	"cmp"
	"context"
	"embed"
	"slices"
	"strings"

	"cogentcore.org/core/colors"
	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/views"
	"github.com/kkoreilly/osusu/osusu"
	"github.com/kkoreilly/osusu/otextencoding"
	"github.com/nlpodyssey/cybertron/pkg/client"
	"github.com/nlpodyssey/cybertron/pkg/models/bert"
	"github.com/nlpodyssey/spago/mat"
	"goki.dev/grows/jsons"
	"goki.dev/grr"
)

//go:embed recipes.json
var recipesFS embed.FS

var recipes []*osusu.Recipe

//go:embed textEncodingVectors.json
var textEncodingVectorsFS embed.FS

var textEncodingVectors map[string][]float32

func configDiscover(rf *core.Frame, mf *core.Frame) {
	if rf.HasChildren() {
		rf.DeleteChildren(true)
	}
	updt := rf.UpdateStart()

	rf.Styler(func(s *styles.Style) {
		s.Wrap = true
	})

	if recipes == nil {
		err := jsons.OpenFS(&recipes, recipesFS, "recipes.json")
		if err != nil {
			core.ErrorDialog(rf, err, "Error opening recipes")
			return
		}
		for _, recipe := range recipes {
			grr.Log(recipe.Init())
		}

		err = jsons.OpenFS(&textEncodingVectors, textEncodingVectorsFS, "textEncodingVectors.json")
		if err != nil {
			core.ErrorDialog(rf, err, "Error opening recipe text encoding vectors")
			return
		}

		otextencoding.Model = client.NewClientForTextEncoding("localhost:8081", client.Options{})
	}

	var meals []*osusu.Meal
	err := osusu.DB.Find(&meals).Error
	if err != nil {
		core.ErrorDialog(rf, err)
	}
	mealEntries := map[uint][]osusu.Entry{}
	mealVectors := make([]mat.Matrix, len(meals))
	for i, meal := range meals {

		entries := []osusu.Entry{}
		err := osusu.DB.Find(&entries, "meal_id = ? AND user_id = ?", meal.ID, curUser.ID).Error
		if err != nil {
			core.ErrorDialog(rf, err)
		}
		mealEntries[meal.ID] = entries

		res, err := otextencoding.Model.Encode(context.TODO(), meal.Text(), int(bert.MeanPooling))
		if err != nil {
			core.ErrorDialog(rf, err, "Error text encoding meal")
			continue
		}
		mealVectors[i] = res.Vector
	}

	for _, recipe := range recipes {
		// first we get the base score index
		// TODO(kai/osusu): cache this step
		recipe.ComputeBaseScoreIndex()

		// then we get the raw text encoding score
		// TODO(kai/osusu): cache this step
		recipeVector := textEncodingVectors[recipe.URL]
		recipeMat := mat.NewDense[float32](mat.WithBacking(recipeVector))
		recipe.TextEncodingScores = map[uint]float32{}
		for i, meal := range meals {
			mealVector := mealVectors[i]
			score := mealVector.DotUnitary(recipeMat)
			recipe.TextEncodingScores[meal.ID] = score.Item().F32()
		}

		// then we get the weighted score
		// this step can not be cached
		weightedScores := make([]*osusu.Score, len(meals))
		for i, meal := range meals {
			textEncodingScore := recipe.TextEncodingScores[meal.ID]
			entries := mealEntries[meal.ID]

			score := meal.Score(entries)
			score.ComputeTotal(curOptions)

			osusu.MulScore(score, textEncodingScore)
			weightedScores[i] = score
		}
		recipe.EncodingScoreIndex = *osusu.AverageScore(weightedScores)
	}

	// now we can compute the percentile scores
	osusu.ComputeNormScores(recipes)

	// and then the total scores
	for _, recipe := range recipes {
		recipe.BaseScore.ComputeTotal(curOptions)
		recipe.EncodingScore.ComputeTotal(curOptions)
		// encoding score is three times more important than base score
		recipe.Score = *osusu.AverageScore([]*osusu.Score{&recipe.BaseScore, &recipe.EncodingScore, &recipe.EncodingScore, &recipe.EncodingScore})
	}

	slices.SortFunc(recipes, func(a, b *osusu.Recipe) int {
		return cmp.Compare(b.Score.Total, a.Score.Total)
	})

	for _, recipe := range recipes {
		recipe := recipe

		if rf.NumChildren() > 100 {
			break
		}

		grr.Log(recipe.CategoryFlag.SetString(strings.Join(recipe.Category, "|")))
		grr.Log(recipe.CuisineFlag.SetString(strings.Join(recipe.Cuisine, "|")))

		if !bitFlagsOverlap(recipe.CategoryFlag, curOptions.Categories) ||
			!bitFlagsOverlap(recipe.CuisineFlag, curOptions.Cuisines) {
			continue
		}

		rc := core.NewFrame(rf)
		cardStyles(rc)

		img := core.NewImage(rc)
		go func() {
			if i := getImageFromURL(recipe.Image); i != nil {
				img.SetImage(i)
				img.Update()
			}
		}()

		core.NewText(rc).SetType(core.TextHeadlineSmall).SetText(recipe.Name)

		castr := friendlyBitFlagString(recipe.CategoryFlag)
		custr := friendlyBitFlagString(recipe.CuisineFlag)
		text := castr
		if castr != "" && custr != "" {
			text += " • "
		}
		text += custr
		core.NewText(rc).SetText(text).Styler(func(s *styles.Style) {
			s.Color = colors.Scheme.OnSurfaceVariant
		})

		scoreGrid(rc, &recipe.Score, true)

		rc.OnClick(func(e events.Event) {
			addRecipe(rf, recipe, rc, mf)
		})
	}

	rf.Update()
	rf.UpdateEndLayout(updt)
}

func addRecipe(rf *core.Frame, recipe *osusu.Recipe, rc *core.Frame, mf *core.Frame) {
	d := core.NewBody().AddTitle("Add recipe")
	views.NewStructView(d).SetStruct(recipe).SetReadOnly(true)
	d.AddBottomBar(func(pw core.Widget) {
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
