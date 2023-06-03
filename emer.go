//go:build !js

package main

import "github.com/emer/emergent/decoder"

// InferCuisinesSoftMax infers the cuisines of the given recipes with cuisines missing using the emergent SoftMax type
func InferCuisinesSoftMax(recipes Recipes) Recipes {
	sm := &decoder.SoftMax{}
	sm.Init(len(allCuisines), len(recipes))
	// for _, recipe := range recipes {
	// 	if len(recipe.Cuisine) != 0 {

	// 	}
	// 	words := GetWords(recipe.Name)

	// 	for i, word := range words {
	// 	}
	// }
	return recipes
}
