package main

import (
	"errors"

	"github.com/kkoreilly/osusu/osusu"
	"goki.dev/colors"
	"goki.dev/cursors"
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
			d := gi.NewDialog(tb).Title("New meal").FullWindow(true)
			meal := &osusu.Meal{}
			giv.NewStructView(d).SetStruct(meal)
			d.OnAccept(func(e events.Event) {
				err := osusu.DB.Create(meal).Error
				if err != nil {
					gi.NewDialog(tb).Title("Error creating meal").Prompt(err.Error()).Ok().Run()
					return
				}
				configMeals(mf)
			}).Cancel().Ok("Create meal").Run()
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
		mc.Style(func(s *styles.Style) {
			s.SetAbilities(true, abilities.Hoverable, abilities.Pressable)
			s.Cursor = cursors.Pointer
			s.Border.Radius = styles.BorderRadiusLarge
			s.BackgroundColor.SetSolid(colors.Scheme.SurfaceContainerLow)
			s.Padding.Set(units.Dp(8))
			s.Min.Set(units.Em(10))
			s.SetGrow(0)
			s.MainAxis = mat32.Y
		})
		gi.NewLabel(mc).SetType(gi.LabelHeadlineSmall).SetText(meal.Name)
		gi.NewLabel(mc).SetText(meal.Description).Style(func(s *styles.Style) {
			s.Color = colors.Scheme.OnSurfaceVariant
		})
		mc.OnClick(func(e events.Event) {
			editMeal(mf, meal, mc)
		})
	}
	mf.Update()
	mf.UpdateEndLayout(updt)
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
