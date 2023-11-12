package main

import (
	"goki.dev/gi/v2/gi"
	"goki.dev/gi/v2/gimain"
)

func main() { gimain.Run(app) }

func app() {
	sc := gi.NewScene("osusu").SetTitle("Osusu")
	gi.NewLabel(sc).SetType(gi.LabelHeadlineSmall).SetText("Osusu")
	gi.NewWindow(sc).Run().Wait()
}
