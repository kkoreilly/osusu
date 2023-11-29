// Package otextencoding provides Osusu text encoding logic.
package otextencoding

import (
	"context"

	"github.com/kkoreilly/osusu/osusu"
	"github.com/nlpodyssey/cybertron/pkg/models/bert"
	"github.com/nlpodyssey/cybertron/pkg/tasks"
	"github.com/nlpodyssey/cybertron/pkg/tasks/textencoding"
)

// Model is the main text encoding model.
var Model textencoding.Interface

// LoadModel loads [Model].
func LoadModel() error {
	m, err := tasks.Load[textencoding.Interface](&tasks.Config{ModelsDir: "models", ModelName: textencoding.DefaultModel})
	if err != nil {
		return err
	}
	Model = m
	return nil
}

// Encode computes the text encoding of the given recipe.
func Encode(recipe *osusu.Recipe) (textencoding.Response, error) {
	return Model.Encode(context.TODO(), recipe.Text(), int(bert.MeanPooling))
}
