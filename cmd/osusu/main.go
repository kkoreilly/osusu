package main

import (
	"cogentcore.org/core/core"
	"github.com/kkoreilly/osusu/osusu"
)

func main() {
	b := core.NewBody("Osusu")
	base(b)
	w := b.NewWindow().Run()
	err := osusu.OpenDB()
	core.ErrorDialog(b, err)
	// loadSession(b)
	w.Wait()
}
