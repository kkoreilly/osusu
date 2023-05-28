package main

// Score is the score of a meal or entry
type Score struct {
	ID          int
	Cost        int
	Effort      int
	Healthiness int
	Taste       int
	Recency     int
	Total       int
}

// ScoreWeight is a combination of a score and a weight for how much it matters
type ScoreWeight struct {
	Score  Score
	Weight int
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

// AverageScore returns a score with the values set to the average values of the given scores, with the given weight given to each of the given scores as specified by the score weights
func AverageScore(scores []ScoreWeight) Score {
	lenScores := len(scores)
	if lenScores == 0 {
		return Score{}
	}
	res := Score{}
	totalWeights := 0
	for _, score := range scores {
		res.Cost += score.Score.Cost * score.Weight
		res.Effort += score.Score.Effort * score.Weight
		res.Healthiness += score.Score.Healthiness * score.Weight
		res.Taste += score.Score.Taste * score.Weight
		res.Recency += score.Score.Recency * score.Weight
		res.Total += score.Score.Total * score.Weight
		totalWeights += score.Weight
	}
	res.Cost /= totalWeights
	res.Effort /= totalWeights
	res.Healthiness /= totalWeights
	res.Taste /= totalWeights
	res.Recency /= totalWeights
	res.Total /= totalWeights
	return res
}
