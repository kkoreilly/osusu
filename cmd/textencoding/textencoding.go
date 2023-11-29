// Command textencoding generates text encoding vectors for all of the database recipes.
package main

import (
	"context"
	"fmt"
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
	grr.Must(jsons.Open(&recipes, "../recipes.json"))

	for _, recipe := range recipes {
		res := grr.Must1(m.Encode(context.TODO(), strings.Join([]string{recipe.Name, recipe.Description}, " "), int(bert.MeanPooling)))
		fmt.Println(res.Vector)
	}
}
