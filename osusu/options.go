package osusu

type Options struct {
	TasteImportance       int `view:"slider" min:"0" def:"50" max:"100"`
	RecencyImportance     int `view:"slider" min:"0" def:"50" max:"100"`
	CostImportance        int `view:"slider" min:"0" def:"50" max:"100"`
	EffortImportance      int `view:"slider" min:"0" def:"50" max:"100"`
	HealthinessImportance int `view:"slider" min:"0" def:"50" max:"100"`
}

func DefaultOptions() *Options {
	return &Options{
		TasteImportance:       50,
		RecencyImportance:     50,
		CostImportance:        50,
		EffortImportance:      50,
		HealthinessImportance: 50,
	}
}
