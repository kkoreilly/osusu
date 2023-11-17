package main

import (
	"github.com/kkoreilly/osusu/osusu"
	"goki.dev/gi/v2/gi"
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
	gi.NewLabel(search).SetType(gi.LabelHeadlineLarge).SetText("Search")
	mf := gi.NewFrame(search)

	var meals []*osusu.Meal
	err := osusu.DB.Find(&meals).Error
	if err != nil {
		gi.NewDialog(mf).Title("Error finding meals").Prompt(err.Error()).Run()
	}
	for _, meal := range meals {
		mc := gi.NewFrame(mf)
		gi.NewLabel(mc).SetText(meal.Name)
	}

	gi.NewWindow(sc).SetSharedWin().Run()
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
