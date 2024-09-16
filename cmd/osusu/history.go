package main

import (
	"cogentcore.org/core/colors"
	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/views"
	"github.com/kkoreilly/osusu/osusu"
)

func configHistory(ef *core.Frame) {
	if ef.HasChildren() {
		ef.DeleteChildren(true)
	}
	updt := ef.UpdateStart()

	entries := []osusu.Entry{}
	err := osusu.DB.Preload("Meal").Find(&entries, "user_id = ?", curUser.ID).Error
	if err != nil {
		core.ErrorDialog(ef, err)
	}
	for _, entry := range entries {
		entry := entry

		if !bitFlagsOverlap(entry.Category, curOptions.Categories) ||
			!bitFlagsOverlap(entry.Source, curOptions.Sources) ||
			!bitFlagsOverlap(entry.Meal.Cuisine, curOptions.Cuisines) {
			continue
		}

		ec := core.NewFrame(ef)
		cardStyles(ec)

		img := core.NewImage(ec)
		go func() {
			if i := getImageFromURL(entry.Meal.Image); i != nil {
				img.SetImage(i)
				img.Update()
			}
		}()

		core.NewText(ec).SetType(core.TextHeadlineSmall).SetText(entry.Time.Format("Monday, January 2, 2006"))

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
		core.NewText(ec).SetText(text).Styler(func(s *styles.Style) {
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

func editEntry(ef *core.Frame, entry *osusu.Entry, ec *core.Frame) {
	d := core.NewBody().AddTitle("Edit entry")
	views.NewStructView(d).SetStruct(entry)
	d.AddBottomBar(func(pw core.Widget) {
		d.AddCancel(pw)
		d.AddOk(pw).SetText("Save").OnClick(func(e events.Event) {
			err := osusu.DB.Save(entry).Error
			if err != nil {
				core.ErrorDialog(d, err)
			}
			configHistory(ef)
		})
	})
	d.NewFullDialog(ec).Run()
}
