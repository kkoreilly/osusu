package osusu

import "cogentcore.org/core/enums"

type Options struct {
	Categories            Categories
	Sources               Sources
	Cuisines              Cuisines
	TasteImportance       int `display:"slider" min:"0" def:"50" max:"100"`
	RecencyImportance     int `display:"slider" min:"0" def:"50" max:"100"`
	CostImportance        int `display:"slider" min:"0" def:"50" max:"100"`
	EffortImportance      int `display:"slider" min:"0" def:"50" max:"100"`
	HealthinessImportance int `display:"slider" min:"0" def:"50" max:"100"`
}

func DefaultOptions() *Options {
	opts := &Options{
		TasteImportance:       50,
		RecencyImportance:     50,
		CostImportance:        50,
		EffortImportance:      50,
		HealthinessImportance: 50,
	}
	for _, v := range opts.Categories.Values() {
		opts.Categories.SetFlag(true, v.(enums.BitFlag))
	}
	for _, v := range opts.Sources.Values() {
		opts.Sources.SetFlag(true, v.(enums.BitFlag))
	}
	for _, v := range opts.Cuisines.Values() {
		opts.Cuisines.SetFlag(true, v.(enums.BitFlag))
	}
	return opts
}
