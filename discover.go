package main

// WordScoreMap returns a map with a score for each word contained in the names of the given meals using the given entries for each meal and options
func WordScoreMap(meals Meals, entriesForEachMeal map[int64]Entries, options Options) map[string]Score {
	sliceMap := map[string][]Score{}
	for _, meal := range meals {
		score := meal.Score(entriesForEachMeal[meal.ID], options)
		words := GetWords(meal.Name)
		for _, word := range words {
			if sliceMap[word] == nil {
				sliceMap[word] = []Score{}
			}
			sliceMap[word] = append(sliceMap[word], score)
		}
	}
	res := map[string]Score{}
	for word, scores := range sliceMap {
		res[word] = AverageScore(scores)
	}
	return res
}
