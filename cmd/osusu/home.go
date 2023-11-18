package main

import (
	"errors"
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
	// gi.NewLabel(sc).SetText("Email address: " + curUser.Email)
	// gi.NewLabel(sc).SetText("Locale: " + curUser.Locale)

	// img := getPicture()
	// if img != nil {
	// 	gi.NewImage(sc).SetImage(img, 0, 0)
	// }

	tabs := gi.NewTabs(sc).SetDeleteTabButtons(false)

	search := tabs.NewTab("Search")

	mf := gi.NewFrame(search)

	configMeals(mf)

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
				configMeals(mf)
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

func configMeals(mf *gi.Frame) {
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
		mc.OnClick(func(e events.Event) {
			gi.NewMenu(func(m *gi.Scene) {
				gi.NewButton(m).SetIcon(icons.Add).SetText("New entry").OnClick(func(e events.Event) {
					newEntry(meal, mc)
				})
				gi.NewButton(m).SetIcon(icons.Visibility).SetText("View entries").OnClick(func(e events.Event) {
					viewEntries(meal, mc)
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

func viewEntries(meal *osusu.Meal, mc *gi.Frame) {
	d := gi.NewDialog(mc).Title("Entries for " + meal.Name).FullWindow(true)
	entries := []osusu.Entry{}
	err := osusu.DB.Find(&entries, "meal_id = ?", meal.ID).Error
	if err != nil {
		gi.NewDialog(d).Title("Error finding entries for meal").Prompt(err.Error()).Ok().Run()
	}
	for _, entry := range entries {
		entry := entry
		ec := gi.NewFrame(d)
		cardStyles(ec)
		gi.NewLabel(ec).SetType(gi.LabelHeadlineSmall).SetText(entry.Time.Format("Monday, January 2, 2006"))
		gi.NewLabel(ec).SetText(friendlyBitFlagString(entry.Category) + " • " + friendlyBitFlagString(entry.Source)).Style(func(s *styles.Style) {
			s.Color = colors.Scheme.OnSurfaceVariant
		})
	}
	d.Run()
}

func editMeal(mf *gi.Frame, meal *osusu.Meal, mc *gi.Frame) {
	d := gi.NewDialog(mc).Title("Edit meal").FullWindow(true)
	giv.NewStructView(d).SetStruct(meal)
	d.OnAccept(func(e events.Event) {
		err := osusu.DB.Save(meal).Error
		if err != nil {
			gi.NewDialog(d).Title("Error saving meal").Prompt(err.Error()).Ok().Run()
		}
		configMeals(mf)
	}).Cancel().Ok("Save").Run()
}

func cardStyles(card *gi.Frame) {
	card.Style(func(s *styles.Style) {
		s.SetAbilities(true, abilities.Hoverable, abilities.Pressable)
		s.Cursor = cursors.Pointer
		s.Border.Radius = styles.BorderRadiusLarge
		s.BackgroundColor.SetSolid(colors.Scheme.SurfaceContainerLow)
		s.Padding.Set(units.Dp(8))
		s.Min.Set(units.Em(10))
		s.SetGrow(0)
		s.MainAxis = mat32.Y
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

// func getPicture() image.Image {
// 	if curUser.Picture == "" {
// 		return nil
// 	}
// 	resp, err := http.Get(curUser.Picture)
// 	if grr.Log0(err) != nil {
// 		return nil
// 	}
// 	defer resp.Body.Close()
// 	img, _, err := images.Read(resp.Body)
// 	if grr.Log0(err) != nil {
// 		return nil
// 	}
// 	return img
// }
