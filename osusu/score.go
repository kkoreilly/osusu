package osusu

import (
	"time"

	"goki.dev/glop/num"
)

type Score struct {
	Taste       int
	Recency     int
	Cost        int
	Effort      int
	Healthiness int
	Total       int
}

func (s *Score) ComputeTotal(opts *Options) {
	s.Total = (s.Taste * opts.TasteImportance) + (s.Recency * opts.RecencyImportance) +
		(s.Cost * opts.CostImportance) + (s.Effort * opts.EffortImportance) +
		(s.Healthiness * opts.HealthinessImportance)

	totImp := opts.TasteImportance + opts.RecencyImportance + opts.CostImportance + opts.EffortImportance + opts.HealthinessImportance
	if totImp != 0 {
		s.Total /= totImp
	}
}

func (e *Entry) Score() *Score {
	s := &Score{
		Taste:       e.Taste,
		Cost:        100 - e.Cost,
		Effort:      100 - e.Effort,
		Healthiness: e.Healthiness,
	}
	return s
}

func (m *Meal) Score(entries []Entry) *Score {
	n := len(entries)
	if n == 0 {
		return &Score{}
	}

	s := &Score{}
	newest := time.Time{}
	for _, entry := range entries {
		s.Taste += entry.Taste
		s.Cost += 100 - entry.Cost
		s.Effort += 100 - entry.Effort
		s.Healthiness += entry.Healthiness
		if entry.Time.After(newest) {
			newest = entry.Time
		}
	}
	s.Taste /= n
	s.Cost /= n
	s.Effort /= n
	s.Healthiness /= n

	days := time.Since(newest) / (24 * time.Hour)
	recency := 2 * days
	s.Recency = min(int(recency), 100)
	return s
}

func MulScore[T num.Number](s *Score, scalar T) {
	s.Taste = int(T(s.Taste) * scalar)
	s.Recency = int(T(s.Recency) * scalar)
	s.Cost = int(T(s.Cost) * scalar)
	s.Effort = int(T(s.Effort) * scalar)
	s.Healthiness = int(T(s.Healthiness) * scalar)
	s.Total = int(T(s.Total) * scalar)
}
