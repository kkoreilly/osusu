package main

import (
	"github.com/kkoreilly/osusu/osusu"
	"goki.dev/colors"
	"goki.dev/gi/v2/gi"
	"goki.dev/gi/v2/giv"
	"goki.dev/girl/styles"
	"goki.dev/goosi/events"
)

func configHistory(ef *gi.Frame) {
	if ef.HasChildren() {
		ef.DeleteChildren(true)
	}
	updt := ef.UpdateStart()

	entries := []osusu.Entry{}
	err := osusu.DB.Preload("Meal").Find(&entries, "user_id = ?", curUser.ID).Error
	if err != nil {
		gi.ErrorDialog(ef, err).Run()
	}
	for _, entry := range entries {
		entry := entry

		if !bitFlagsOverlap(entry.Category, curOptions.Categories) ||
			!bitFlagsOverlap(entry.Source, curOptions.Sources) ||
			!bitFlagsOverlap(entry.Meal.Cuisine, curOptions.Cuisines) {
			continue
		}

		ec := gi.NewFrame(ef)
		cardStyles(ec)

		img := getImageFromURL(entry.Meal.Image)
		if img != nil {
			gi.NewImage(ec).SetImage(img)
		}

		gi.NewLabel(ec).SetType(gi.LabelHeadlineSmall).SetText(entry.Time.Format("Monday, January 2, 2006"))

		castr := friendlyBitFlagString(entry.Category)
		sostr := friendlyBitFlagString(entry.Source)
		text := entry.Meal.Name
		if entry.Meal.Name != "" && castr != "" {
			text += " • "
		}
		text += castr
		if castr != "" && sostr != "" {
			text += " • "
		}
		text += sostr
		gi.NewLabel(ec).SetText(text).Style(func(s *styles.Style) {
			s.Color = colors.Scheme.OnSurfaceVariant
		})

		score := entry.Score()
		score.ComputeTotal(curOptions)
		scoreGrid(ec, score, false)

		ec.OnClick(func(e events.Event) {
			editEntry(ef, &entry, ec)
		})
	}

	ef.Update()
	ef.UpdateEndLayout(updt)
}

func editEntry(ef *gi.Frame, entry *osusu.Entry, ec *gi.Frame) {
	d := gi.NewBody().AddTitle("Edit entry")
	giv.NewStructView(d).SetStruct(entry)
	d.AddBottomBar(func(pw gi.Widget) {
		d.AddCancel(pw)
		d.AddOk(pw).SetText("Save").OnClick(func(e events.Event) {
			err := osusu.DB.Save(entry).Error
			if err != nil {
				gi.ErrorDialog(d, err).Run()
			}
			configHistory(ef)
		})
	})
	d.NewFullDialog(ec).Run()
}
