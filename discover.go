package main

// WordScoreMap returns a map with a score for each word contained in the names of the given meals using the given entries for each meal and options
func WordScoreMap(meals Meals, entriesForEachMeal map[int64]Entries, options Options) map[string]Score {
	sliceMap := map[string][]ScoreWeight{}
	for _, meal := range meals {
		score := meal.Score(entriesForEachMeal[meal.ID], options)
		words := GetWords(meal.Name)
		for i, word := range words {
			if sliceMap[word] == nil {
				sliceMap[word] = []ScoreWeight{}
			}
			// earlier words are more important
			sliceMap[word] = append(sliceMap[word], ScoreWeight{Score: score, Weight: len(words) - i})
		}
	}
	res := map[string]Score{}
	for word, scores := range sliceMap {
		res[word] = AverageScore(scores)
	}
	return res
}
