package main

import (
	"time"

	"cogentcore.org/core/colors"
	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/styles"
	"github.com/kkoreilly/osusu/osusu"
)

func configSearch(mf *core.Frame) {
	if mf.HasChildren() {
		mf.DeleteChildren(true)
	}
	updt := mf.UpdateStart()

	mf.Styler(func(s *styles.Style) {
		s.Wrap = true
	})

	var meals []*osusu.Meal
	err := osusu.DB.Find(&meals).Error
	if err != nil {
		core.ErrorDialog(mf, err)
	}
	for _, meal := range meals {
		meal := meal

		if !bitFlagsOverlap(meal.Category, curOptions.Categories) ||
			!bitFlagsOverlap(meal.Source, curOptions.Sources) ||
			!bitFlagsOverlap(meal.Cuisine, curOptions.Cuisines) {
			continue
		}

		mc := core.NewFrame(mf)
		cardStyles(mc)

		img := core.NewImage(mc)
		go func() {
			if i := getImageFromURL(meal.Image); i != nil {
				img.SetImage(i)
				img.Update()
			}
		}()

		core.NewText(mc).SetType(core.TextHeadlineSmall).SetText(meal.Name)

		castr := friendlyBitFlagString(meal.Category)
		custr := friendlyBitFlagString(meal.Cuisine)
		text := castr
		if castr != "" && custr != "" {
			text += " • "
		}
		text += custr
		core.NewText(mc).SetText(text).Styler(func(s *styles.Style) {
			s.Color = colors.Scheme.OnSurfaceVariant
		})

		entries := []osusu.Entry{}
		err := osusu.DB.Find(&entries, "meal_id = ? AND user_id = ?", meal.ID, curUser.ID).Error
		if err != nil {
			core.ErrorDialog(mc, err)
		}

		score := meal.Score(entries)
		score.ComputeTotal(curOptions)
		scoreGrid(mc, score, true)

		mc.OnClick(func(e events.Event) {
			core.NewMenu(func(m *core.Scene) {
				core.NewButton(m).SetIcon(icons.Add).SetText("New entry").OnClick(func(e events.Event) {
					newEntry(meal, mc)
				})
				core.NewButton(m).SetIcon(icons.Visibility).SetText("View entries").OnClick(func(e events.Event) {
					viewEntries(meal, entries, mc)
				})
				core.NewButton(m).SetIcon(icons.Edit).SetText("Edit meal").OnClick(func(e events.Event) {
					editMeal(mf, meal, mc)
				})
			}, mc, mc.ContextMenuPos(e)).Run()
		})
	}
	mf.Update()
	mf.UpdateEndLayout(updt)
}

func newEntry(meal *osusu.Meal, mc *core.Frame) {
	d := core.NewBody().AddTitle("Create entry")
	entry := &osusu.Entry{
		MealID:      meal.ID,
		UserID:      curUser.ID,
		Time:        time.Now(),
		Cost:        50,
		Effort:      50,
		Healthiness: 50,
		Taste:       50,
	}
	core.NewForm(d).SetStruct(entry)
	d.AddBottomBar(func(pw core.Widget) {
		d.AddCancel(pw)
		d.AddOk(pw).SetText("Create").OnClick(func(e events.Event) {
			err := osusu.DB.Create(entry).Error
			if err != nil {
				core.ErrorDialog(d, err)
			}
		})
	})
	d.NewFullDialog(mc).Run()
}

func viewEntries(meal *osusu.Meal, entries []osusu.Entry, mc *core.Frame) {
	d := core.NewBody().AddTitle("Entries for " + meal.Name)
	d.AddTopBar(func(pw core.Widget) {
		tb := d.DefaultTopAppBar(pw)
		core.NewButton(tb).SetIcon(icons.Add).SetText("New entry").OnClick(func(e events.Event) {
			newEntry(meal, mc)
		})
	})
	for _, entry := range entries {
		entry := &entry
		ec := core.NewFrame(d)
		cardStyles(ec)
		core.NewText(ec).SetType(core.TextHeadlineSmall).SetText(entry.Time.Format("Monday, January 2, 2006"))

		castr := friendlyBitFlagString(entry.Category)
		sostr := friendlyBitFlagString(entry.Source)
		text := castr
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
			d := core.NewBody().AddTitle("Edit entry")
			core.NewForm(d).SetStruct(entry)
			d.AddBottomBar(func(pw core.Widget) {
				d.AddCancel(pw)
				d.AddOk(pw).SetText("Save").OnClick(func(e events.Event) {
					err := osusu.DB.Save(entry).Error
					if err != nil {
						core.ErrorDialog(d, err)
					}
				})
			})
			d.NewFullDialog(ec).Run()
		})
	}
	d.NewFullDialog(mc).Run()
}

func editMeal(mf *core.Frame, meal *osusu.Meal, mc *core.Frame) {
	d := core.NewBody().AddTitle("Edit meal")
	core.NewForm(d).SetStruct(meal)
	d.AddBottomBar(func(pw core.Widget) {
		d.AddCancel(pw)
		d.AddOk(pw).SetText("Save").OnClick(func(e events.Event) {
			err := osusu.DB.Save(meal).Error
			if err != nil {
				core.ErrorDialog(d, err)
			}
			configSearch(mf)
		})
	})
	d.NewFullDialog(mc).Run()
}
