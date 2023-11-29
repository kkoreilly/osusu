package main

import (
	"errors"
	"image"
	"net/http"
	"strconv"

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
	"goki.dev/grows/images"
	"goki.dev/grr"
	"goki.dev/icons"
	"gorm.io/gorm"
)

var (
	curUser    *osusu.User
	curGroup   *osusu.Group
	curOptions = osusu.DefaultOptions()
)

func home() {
	b := gi.NewBody("home")

	tabs := gi.NewTabs(b).SetDeleteTabButtons(false)

	search := tabs.NewTab("Search")
	mf := gi.NewFrame(search)
	configSearch(mf)

	discover := tabs.NewTab("Discover")
	rf := gi.NewFrame(discover)
	tabs.Tabs().ChildByName("discover").(gi.Widget).OnClick(func(e events.Event) {
		if !rf.HasChildren() {
			configDiscover(rf, mf)
		}
	})

	history := tabs.NewTab("History")
	ef := gi.NewFrame(history)
	configHistory(ef)

	b.AddTopBar(func(pw gi.Widget) {
		tb := b.DefaultTopAppBar(pw)
		gi.NewButton(tb).SetIcon(icons.Add).SetText("New meal").OnClick(func(e events.Event) {
			newMeal(tb, mf, &osusu.Meal{})
		})
		gi.NewButton(tb).SetIcon(icons.Sort).SetText("Sort").OnClick(func(e events.Event) {
			d := gi.NewBody().AddTitle("Sort and filter")
			giv.NewStructView(d).SetStruct(curOptions)
			d.OnClose(func(e events.Event) {
				configSearch(mf)
				configHistory(ef)
				configDiscover(mf, rf)
			})
			d.NewFullDialog(tb).Run()
		})
	})

	b.NewWindow().Run()

	curGroup = &osusu.Group{}
	err := osusu.DB.First(curGroup, curUser.GroupID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			groups(b)
		} else {
			gi.ErrorDialog(b, err)
		}
	}
}

func newMeal(ctx gi.Widget, mf *gi.Frame, meal *osusu.Meal) {
	d := gi.NewBody().AddTitle("Create meal")
	giv.NewStructView(d).SetStruct(meal)
	d.AddBottomBar(func(pw gi.Widget) {
		d.AddCancel(pw)
		d.AddOk(pw).SetText("Create").OnClick(func(e events.Event) {
			err := osusu.DB.Create(meal).Error
			if err != nil {
				gi.ErrorDialog(d, err)
				return
			}
			configSearch(mf)
		})
	})
	d.NewFullDialog(ctx).Run()
}

func cardStyles(card *gi.Frame) {
	card.Style(func(s *styles.Style) {
		s.SetAbilities(true, abilities.Hoverable, abilities.Pressable)
		s.Cursor = cursors.Pointer
		s.Border.Radius = styles.BorderRadiusLarge
		s.BackgroundColor.SetSolid(colors.Scheme.SurfaceContainerLow)
		s.Padding.Set(units.Dp(8))
		s.Min.Set(units.Em(20))
		s.Grow.Set(0, 0)
		s.Direction = styles.Column
		s.Align.Content = styles.Center
		s.Align.Items = styles.Center
	})
	card.OnWidgetAdded(func(w gi.Widget) {
		switch w := w.(type) {
		case *gi.Label:
			w.Style(func(s *styles.Style) {
				s.SetNonSelectable()
				s.Text.Align = styles.Center
			})
		case *gi.Image:
			w.Style(func(s *styles.Style) {
				s.Min.Set(units.Em(20))
				s.ObjectFit = styles.FitCover
			})
		}
	})
}

func scoreGrid(card *gi.Frame, score *osusu.Score, showRecency bool) *gi.Layout {
	grid := gi.NewLayout(card)
	grid.Style(func(s *styles.Style) {
		s.Display = styles.Grid
		if showRecency {
			s.Columns = 6
		} else {
			s.Columns = 5
		}
		s.Justify.Content = styles.Center
		s.Justify.Items = styles.Center
	})

	label := func(text string) {
		gi.NewLabel(grid).SetType(gi.LabelLabelLarge).SetText(text)
	}

	label("Total")
	label("Taste")
	if showRecency {
		label("New")
	}
	label("Cost")
	label("Effort")
	label("Health")

	value := func(value int) {
		gi.NewLabel(grid).SetText(strconv.Itoa(value))
	}

	value(score.Total)
	value(score.Taste)
	if showRecency {
		value(score.Recency)
	}
	value(score.Cost)
	value(score.Effort)
	value(score.Healthiness)
	return grid
}

func getImageFromURL(url string) image.Image {
	if url == "" {
		return nil
	}
	resp, err := http.Get(url)
	if grr.Log(err) != nil {
		return nil
	}
	defer resp.Body.Close()
	img, _, err := images.Read(resp.Body)
	if grr.Log(err) != nil {
		return nil
	}
	return img
}

// bitFlagsOverlap returns whether there is any overlap between the two bit flags.
// They should be of the same type.
func bitFlagsOverlap(a, b enums.BitFlag) bool {
	vals := a.Values()
	for _, v := range vals {
		vb := v.(enums.BitFlag)
		if a.HasFlag(vb) && b.HasFlag(vb) {
			return true
		}
	}
	return false
}

func friendlyBitFlagString(bf enums.BitFlag) string {
	matches := []string{}
	vals := bf.Values()
	for _, v := range vals {
		vb := v.(enums.BitFlag)
		if bf.HasFlag(vb) {
			vbs := vb.BitIndexString()
			matches = append(matches, vbs)
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
