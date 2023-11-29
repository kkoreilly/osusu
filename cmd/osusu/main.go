package main

import (
	"github.com/kkoreilly/osusu/osusu"
	"goki.dev/gi/v2/gi"
	"goki.dev/gi/v2/gimain"
)

func main() { gimain.Run(app) }

func app() {
	gi.SetAppName("osusu")
	b := gi.NewBody().SetTitle("Osusu")
	base(b)
	w := b.NewWindow().Run()
	err := osusu.OpenDB()
	if err != nil {
		gi.ErrorDialog(b, err)
	}
	// loadSession(b)
	w.Wait()
}
