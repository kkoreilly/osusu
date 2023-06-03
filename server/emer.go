//go:build !js

package server

import (
	"github.com/emer/emergent/decoder"
	"github.com/kkoreilly/osusu/osusu"
)

// InferCuisinesSoftMax infers the cuisines of the given recipes with cuisines missing using the emergent SoftMax type
func InferCuisinesSoftMax(recipes osusu.Recipes) osusu.Recipes {
	sm := decoder.SoftMax{}
	sm.Init(len(osusu.AllCuisines), len(recipes))

	// for _, recipe := range recipes {
	// 	if len(recipe.Cuisine) != 0 {

	// 	}
	// 	words := GetWords(recipe.Name)

	// 	for i, word := range words {
	// 	}
	// }
	return recipes
}
