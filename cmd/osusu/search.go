package main

import (
	"time"

	"github.com/kkoreilly/osusu/osusu"
	"goki.dev/colors"
	"goki.dev/gi/v2/gi"
	"goki.dev/gi/v2/giv"
	"goki.dev/girl/styles"
	"goki.dev/goosi/events"
	"goki.dev/icons"
)

func configSearch(mf *gi.Frame) {
	if mf.HasChildren() {
		mf.DeleteChildren(true)
	}
	updt := mf.UpdateStart()

	mf.Style(func(s *styles.Style) {
		s.Wrap = true
	})

	var meals []*osusu.Meal
	err := osusu.DB.Find(&meals).Error
	if err != nil {
		gi.ErrorDialog(mf, err).Run()
	}
	for _, meal := range meals {
		meal := meal

		if !bitFlagsOverlap(meal.Category, curOptions.Categories) ||
			!bitFlagsOverlap(meal.Source, curOptions.Sources) ||
			!bitFlagsOverlap(meal.Cuisine, curOptions.Cuisines) {
			continue
		}

		mc := gi.NewFrame(mf)
		cardStyles(mc)

		if meal.Image != "" {
			img := gi.NewImage(mc)
			go func() {
				img.SetImage(getImageFromURL(meal.Image))
				img.Update()
			}()
		}

		gi.NewLabel(mc).SetType(gi.LabelHeadlineSmall).SetText(meal.Name)

		castr := friendlyBitFlagString(meal.Category)
		custr := friendlyBitFlagString(meal.Cuisine)
		text := castr
		if castr != "" && custr != "" {
			text += " • "
		}
		text += custr
		gi.NewLabel(mc).SetText(text).Style(func(s *styles.Style) {
			s.Color = colors.Scheme.OnSurfaceVariant
		})

		entries := []osusu.Entry{}
		err := osusu.DB.Find(&entries, "meal_id = ? AND user_id = ?", meal.ID, curUser.ID).Error
		if err != nil {
			gi.ErrorDialog(mc, err).Run()
		}

		score := meal.Score(entries)
		score.ComputeTotal(curOptions)
		scoreGrid(mc, score, true)

		mc.OnClick(func(e events.Event) {
			gi.NewMenu(func(m *gi.Scene) {
				gi.NewButton(m).SetIcon(icons.Add).SetText("New entry").OnClick(func(e events.Event) {
					newEntry(meal, mc)
				})
				gi.NewButton(m).SetIcon(icons.Visibility).SetText("View entries").OnClick(func(e events.Event) {
					viewEntries(meal, entries, mc)
				})
				gi.NewButton(m).SetIcon(icons.Edit).SetText("Edit meal").OnClick(func(e events.Event) {
					editMeal(mf, meal, mc)
				})
			}, mc, mc.ContextMenuPos(e)).Run()
		})
	}
	mf.Update()
	mf.UpdateEndLayout(updt)
}

func newEntry(meal *osusu.Meal, mc *gi.Frame) {
	d := gi.NewBody().AddTitle("Create entry")
	entry := &osusu.Entry{
		MealID:      meal.ID,
		UserID:      curUser.ID,
		Time:        time.Now(),
		Cost:        50,
		Effort:      50,
		Healthiness: 50,
		Taste:       50,
	}
	giv.NewStructView(d).SetStruct(entry)
	d.AddBottomBar(func(pw gi.Widget) {
		d.AddCancel(pw)
		d.AddOk(pw).SetText("Create").OnClick(func(e events.Event) {
			err := osusu.DB.Create(entry).Error
			if err != nil {
				gi.ErrorDialog(d, err).Run()
			}
		})
	})
}

func viewEntries(meal *osusu.Meal, entries []osusu.Entry, mc *gi.Frame) {
	d := gi.NewBody().AddTitle("Entries for " + meal.Name)
	d.AddTopBar(func(pw gi.Widget) {
		tb := d.DefaultTopAppBar(pw)
		gi.NewButton(tb).SetIcon(icons.Add).SetText("New entry").OnClick(func(e events.Event) {
			newEntry(meal, mc)
		})
	})
	for _, entry := range entries {
		entry := &entry
		ec := gi.NewFrame(d)
		cardStyles(ec)
		gi.NewLabel(ec).SetType(gi.LabelHeadlineSmall).SetText(entry.Time.Format("Monday, January 2, 2006"))

		castr := friendlyBitFlagString(entry.Category)
		sostr := friendlyBitFlagString(entry.Source)
		text := castr
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
			d := gi.NewBody().AddTitle("Edit entry")
			giv.NewStructView(d).SetStruct(entry)
			d.AddBottomBar(func(pw gi.Widget) {
				d.AddCancel(pw)
				d.AddOk(pw).SetText("Save").OnClick(func(e events.Event) {
					err := osusu.DB.Save(entry).Error
					if err != nil {
						gi.ErrorDialog(d, err).Run()
					}
				})
			})
			d.NewFullDialog(ec).Run()
		})
	}
	d.NewFullDialog(mc).Run()
}

func editMeal(mf *gi.Frame, meal *osusu.Meal, mc *gi.Frame) {
	d := gi.NewBody().AddTitle("Edit meal")
	giv.NewStructView(d).SetStruct(meal)
	d.AddBottomBar(func(pw gi.Widget) {
		d.AddCancel(pw)
		d.AddOk(pw).SetText("Save").OnClick(func(e events.Event) {
			err := osusu.DB.Save(meal).Error
			if err != nil {
				gi.ErrorDialog(d, err).Run()
			}
			configSearch(mf)
		})
	})
	d.NewFullDialog(mc).Run()
}
