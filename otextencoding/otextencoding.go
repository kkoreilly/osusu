// Package otextencoding provides Osusu text encoding logic.
package otextencoding

import (
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
