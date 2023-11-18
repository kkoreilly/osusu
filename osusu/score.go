package osusu

import (
	"time"
)

type Score struct {
	Taste       int
	Cost        int
	Effort      int
	Healthiness int
	Recency     int
	Total       int
}

func (s *Score) ComputeTotal() {
	s.Total = s.Taste + s.Cost + s.Effort + s.Healthiness + s.Recency
	s.Total /= 5
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