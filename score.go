package main

// Score is the score of a meal or entry
type Score struct {
	Cost        int
	Effort      int
	Healthiness int
	Taste       int
	Recency     int
	Total       int
}

// ComputeTotal computes and returns the total score of the score.
func (s Score) ComputeTotal(options Options) int {
	den := options.CostWeight + options.EffortWeight + options.HealthinessWeight + options.TasteWeight + options.RecencyWeight
	if den == 0 {
		return 0
	}
	sum := s.Cost*options.CostWeight + s.Effort*options.EffortWeight + s.Healthiness*options.HealthinessWeight + s.Taste*options.TasteWeight + s.Recency*options.RecencyWeight
	return sum / den
}

// AverageScore returns a score with the values set to the average values of the given scores
func AverageScore(scores []Score) Score {
	lenScores := len(scores)
	if lenScores == 0 {
		return Score{}
	}
	res := Score{}
	for _, score := range scores {
		res.Cost += score.Cost
		res.Effort += score.Effort
		res.Healthiness += score.Healthiness
		res.Taste += score.Taste
		res.Recency += score.Recency
		res.Total += score.Total
	}
	res.Cost /= lenScores
	res.Effort /= lenScores
	res.Healthiness /= lenScores
	res.Taste /= lenScores
	res.Recency /= lenScores
	res.Total /= lenScores
	return res
}
