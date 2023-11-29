// Command textencoding generates text encoding vectors for all of the database recipes.
package main

import (
	"context"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/kkoreilly/osusu/osusu"
	"github.com/kkoreilly/osusu/otextencoding"
	"github.com/nlpodyssey/cybertron/pkg/models/bert"
	"github.com/rs/zerolog"
	"goki.dev/grows/jsons"
	"goki.dev/grr"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	grr.Must(otextencoding.LoadModel())

	var recipes []*osusu.Recipe
	grr.Must(jsons.Open(&recipes, filepath.Join("..", "osusu", "recipes.json")))

	// keyed by recipe URL
	vectors := map[string][]float32{}

	st := time.Now()
	nrecipes := len(recipes)

	slog.Info("starting", "numRecipes", nrecipes)

	for i, recipe := range recipes {
		res := grr.Must1(otextencoding.Model.Encode(context.TODO(), recipe.Text(), int(bert.MeanPooling)))
		vectors[recipe.URL] = res.Vector.Data().F32()
		if i%10 == 0 && i != 0 {
			slog.Info("on", "recipe", i, "estimatedTimeRemaining", time.Since(st)*time.Duration((nrecipes-i)/i))
		}
	}

	grr.Must(jsons.Save(vectors, filepath.Join("..", "osusu", "textEncodingVectors.json")))
}
