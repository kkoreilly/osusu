package main

import (
	"github.com/kkoreilly/osusu/osusu"
	"goki.dev/colors"
	"goki.dev/gi/v2/gi"
	"goki.dev/gi/v2/giv"
	"goki.dev/girl/styles"
	"goki.dev/girl/units"
	"goki.dev/goosi/events"
	"goki.dev/icons"
)

var curUser *osusu.User

func home() {
	sc := gi.NewScene("home")
	gi.NewLabel(sc).SetType(gi.LabelHeadlineLarge).SetText("Welcome, " + curUser.Name)
	// gi.NewLabel(sc).SetText("Email address: " + curUser.Email)
	// gi.NewLabel(sc).SetText("Locale: " + curUser.Locale)

	// img := getPicture()
	// if img != nil {
	// 	gi.NewImage(sc).SetImage(img, 0, 0)
	// }

	tabs := gi.NewTabs(sc)

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
					gi.NewDialog(d).Title("Error creating meal").Prompt(err.Error()).Run()
				}
				configMeals(mf)
			}).Cancel().Ok().Run()
		})
	}

	gi.NewWindow(sc).SetSharedWin().Run()
}

func configMeals(mf *gi.Frame) {
	if mf.HasChildren() {
		mf.DeleteChildren(true)
	}
	updt := mf.UpdateStart()

	var meals []*osusu.Meal
	err := osusu.DB.Find(&meals).Error
	if err != nil {
		gi.NewDialog(mf).Title("Error finding meals").Prompt(err.Error()).Run()
	}
	for _, meal := range meals {
		mc := gi.NewFrame(mf).Style(func(s *styles.Style) {
			s.Border.Radius = styles.BorderRadiusLarge
			s.BackgroundColor.SetSolid(colors.Scheme.SurfaceContainerLow)
			s.Padding.Set(units.Dp(8))
			s.Min.Set(units.Em(10))
			s.SetGrow(0)
		})
		gi.NewLabel(mc).SetType(gi.LabelHeadlineSmall).SetText(meal.Name)
	}
	mf.Update()
	mf.UpdateEndLayout(updt)
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
