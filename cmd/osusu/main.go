package main

import (
	"github.com/kkoreilly/osusu/osusu"
	"goki.dev/gi/v2/gi"
	"goki.dev/gi/v2/gimain"
)

func main() { gimain.Run(app) }

func app() {
	sc := gi.NewScene("osusu").SetTitle("Osusu")
	base(sc)
	w := gi.NewWindow(sc).Run()
	err := osusu.OpenDB()
	if err != nil {
		gi.NewDialog(sc).Title("Error opening database").Prompt(err.Error()).Run()
	}
	w.Wait()
}
