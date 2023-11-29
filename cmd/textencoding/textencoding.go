// Command textencoding generates text encoding vectors for all of the database recipes.
package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/kkoreilly/osusu/osusu"
	"github.com/nlpodyssey/cybertron/pkg/models/bert"
	"github.com/nlpodyssey/cybertron/pkg/tasks"
	"github.com/nlpodyssey/cybertron/pkg/tasks/textencoding"
	"github.com/rs/zerolog"
	"goki.dev/grows/jsons"
	"goki.dev/grr"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	m := grr.Must1(tasks.Load[textencoding.Interface](&tasks.Config{ModelsDir: "models", ModelName: textencoding.DefaultModel}))

	var recipes []osusu.Recipe
	grr.Must(jsons.Open(&recipes, filepath.Join("..", "osusu", "recipes.json")))

	// keyed by recipe URL
	vectors := map[string][]float32{}

	for i, recipe := range recipes {
		rstr := strings.Join([]string{recipe.Name, recipe.Description, strings.Join(recipe.Ingredients, " ")}, " ")
		res := grr.Must1(m.Encode(context.TODO(), rstr, int(bert.MeanPooling)))
		vectors[recipe.URL] = res.Vector.Data().F32()
		if i%10 == 0 {
			fmt.Println("On recipe", i)
		}
		if i == 100 {
			break
		}
	}

	grr.Must(jsons.Save(vectors, "textencodingvectors.json"))
}
