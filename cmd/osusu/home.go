package main

import (
	"errors"
	"strconv"
	"time"

	"github.com/kkoreilly/osusu/osusu"
	"goki.dev/colors"
	"goki.dev/cursors"
	"goki.dev/enums"
	"goki.dev/gi/v2/gi"
	"goki.dev/gi/v2/giv"
	"goki.dev/girl/abilities"
	"goki.dev/girl/styles"
	"goki.dev/girl/units"
	"goki.dev/goosi/events"
	"goki.dev/icons"
	"goki.dev/mat32/v2"
	"gorm.io/gorm"
)

var curUser *osusu.User
var curGroup *osusu.Group

func home() {
	sc := gi.NewScene("home")
	gi.NewLabel(sc).SetType(gi.LabelHeadlineLarge).SetText("Welcome, " + curUser.Name)

	tabs := gi.NewTabs(sc).SetDeleteTabButtons(false)

	search := tabs.NewTab("Search")
	mf := gi.NewFrame(search)
	configSearch(mf)

	history := tabs.NewTab("History")
	ef := gi.NewFrame(history)
	configHistory(ef)

	gi.DefaultTopAppBar = func(tb *gi.TopAppBar) {
		gi.DefaultTopAppBarStd(tb)
		gi.NewButton(tb).SetIcon(icons.Add).SetText("New meal").OnClick(func(e events.Event) {
			d := gi.NewDialog(tb).Title("Create meal").FullWindow(true)
			meal := &osusu.Meal{}
			giv.NewStructView(d).SetStruct(meal)
			d.OnAccept(func(e events.Event) {
				err := osusu.DB.Create(meal).Error
				if err != nil {
					gi.NewDialog(tb).Title("Error creating meal").Prompt(err.Error()).Ok().Run()
					return
				}
				configSearch(mf)
			}).Cancel().Ok("Create").Run()
		})
	}

	gi.NewWindow(sc).SetNewWindow(false).Run()

	curGroup = &osusu.Group{}
	err := osusu.DB.First(curGroup, curUser.GroupID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			groups(sc)
		} else {
			gi.NewDialog(sc).Title("Error getting current group").Prompt(err.Error()).Run()
		}
	}
}

func configSearch(mf *gi.Frame) {
	if mf.HasChildren() {
		mf.DeleteChildren(true)
	}
	updt := mf.UpdateStart()

	var meals []*osusu.Meal
	err := osusu.DB.Find(&meals).Error
	if err != nil {
		gi.NewDialog(mf).Title("Error finding meals").Prompt(err.Error()).Ok().Run()
	}
	for _, meal := range meals {
		meal := meal

		mc := gi.NewFrame(mf)
		cardStyles(mc)
		gi.NewLabel(mc).SetType(gi.LabelHeadlineSmall).SetText(meal.Name)
		gi.NewLabel(mc).SetText(friendlyBitFlagString(meal.Category) + " • " + friendlyBitFlagString(meal.Cuisine)).Style(func(s *styles.Style) {
			s.Color = colors.Scheme.OnSurfaceVariant
		})

		entries := []osusu.Entry{}
		err := osusu.DB.Find(&entries, "meal_id = ? AND user_id = ?", meal.ID, curUser.ID).Error
		if err != nil {
			gi.NewDialog(mc).Title("Error finding entries for meal").Prompt(err.Error()).Ok().Run()
		}

		score := meal.Score(entries)
		score.ComputeTotal()

		gi.NewLabel(mc).SetText("Score: " + strconv.Itoa(score.Total))

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
	d := gi.NewDialog(mc).Title("Create entry").FullWindow(true)
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
	d.OnAccept(func(e events.Event) {
		err := osusu.DB.Create(entry).Error
		if err != nil {
			gi.NewDialog(d).Title("Error creating entry").Prompt(err.Error()).Ok().Run()
		}
	}).Cancel().Ok("Create").Run()
}

func viewEntries(meal *osusu.Meal, entries []osusu.Entry, mc *gi.Frame) {
	d := gi.NewDialog(mc).Title("Entries for " + meal.Name).FullWindow(true)
	d.TopAppBar = func(tb *gi.TopAppBar) {
		gi.DefaultTopAppBarStd(tb)
		gi.NewButton(tb).SetIcon(icons.Add).SetText("New entry").OnClick(func(e events.Event) {
			newEntry(meal, mc)
		})
	}
	for _, entry := range entries {
		entry := &entry
		ec := gi.NewFrame(d)
		cardStyles(ec)
		gi.NewLabel(ec).SetType(gi.LabelHeadlineSmall).SetText(entry.Time.Format("Monday, January 2, 2006"))
		gi.NewLabel(ec).SetText(friendlyBitFlagString(entry.Category) + " • " + friendlyBitFlagString(entry.Source)).Style(func(s *styles.Style) {
			s.Color = colors.Scheme.OnSurfaceVariant
		})
		ec.OnClick(func(e events.Event) {
			editEntry(entry, ec)
		})
	}
	d.Run()
}

func editEntry(entry *osusu.Entry, ec *gi.Frame) {
	d := gi.NewDialog(ec).Title("Edit entry").FullWindow(true)
	giv.NewStructView(d).SetStruct(entry)
	d.OnAccept(func(e events.Event) {
		err := osusu.DB.Save(entry).Error
		if err != nil {
			gi.NewDialog(d).Title("Error saving entry").Prompt(err.Error()).Ok().Run()
		}
	}).Cancel().Ok("Save").Run()
}

func editMeal(mf *gi.Frame, meal *osusu.Meal, mc *gi.Frame) {
	d := gi.NewDialog(mc).Title("Edit meal").FullWindow(true)
	giv.NewStructView(d).SetStruct(meal)
	d.OnAccept(func(e events.Event) {
		err := osusu.DB.Save(meal).Error
		if err != nil {
			gi.NewDialog(d).Title("Error saving meal").Prompt(err.Error()).Ok().Run()
		}
		configSearch(mf)
	}).Cancel().Ok("Save").Run()
}

func configHistory(ef *gi.Frame) {
	entries := []osusu.Entry{}
	err := osusu.DB.Preload("Meal").Find(&entries, "user_id = ?", curUser.ID).Error
	if err != nil {
		gi.NewDialog(ef).Title("Error finding entries").Prompt(err.Error()).Ok().Run()
	}
	for _, entry := range entries {
		entry := &entry
		ec := gi.NewFrame(ef)
		cardStyles(ec)
		gi.NewLabel(ec).SetType(gi.LabelHeadlineSmall).SetText(entry.Time.Format("Monday, January 2, 2006"))
		gi.NewLabel(ec).SetText(entry.Meal.Name + " • " + friendlyBitFlagString(entry.Category) + " • " + friendlyBitFlagString(entry.Source)).Style(func(s *styles.Style) {
			s.Color = colors.Scheme.OnSurfaceVariant
		})
		ec.OnClick(func(e events.Event) {
			editEntry(entry, ec)
		})
	}
}

func cardStyles(card *gi.Frame) {
	card.Style(func(s *styles.Style) {
		s.SetAbilities(true, abilities.Hoverable, abilities.Pressable)
		s.Cursor = cursors.Pointer
		s.Border.Radius = styles.BorderRadiusLarge
		s.BackgroundColor.SetSolid(colors.Scheme.SurfaceContainerLow)
		s.Padding.Set(units.Dp(8))
		s.Min.Set(units.Em(20))
		s.SetGrow(0)
		s.MainAxis = mat32.Y
	})
	card.OnWidgetAdded(func(w gi.Widget) {
		if _, ok := w.(*gi.Label); ok {
			w.Style(func(s *styles.Style) {
				s.SetNonSelectable()
			})
		}
	})
}

func friendlyBitFlagString(bf enums.BitFlag) string {
	matches := []string{}
	vals := bf.Values()
	for _, e := range vals {
		eb := e.(enums.BitFlag)
		if bf.HasFlag(eb) {
			ebs := eb.BitIndexString()
			matches = append(matches, ebs)
		}
	}
	switch len(matches) {
	case 0:
		return ""
	case 1:
		return matches[0]
	case 2:
		return matches[0] + " and " + matches[1]
	}
	res := ""
	for i, match := range matches {
		res += match
		if i == len(matches)-1 {
			// last one, so do nothing
		} else if i == len(matches)-2 {
			res += ", and "
		} else {
			res += ", "
		}
	}
	return res
}
