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
	"goki.dev/mat32/v2"
	"gorm.io/gorm"
)

var (
	curUser    *osusu.User
	curGroup   *osusu.Group
	curOptions = osusu.DefaultOptions()
)

func home() {
	sc := gi.NewScene("home")

	tabs := gi.NewTabs(sc).SetDeleteTabButtons(false)

	search := tabs.NewTab("Search")
	mf := gi.NewFrame(search)
	configSearch(mf)

	discover := tabs.NewTab("Discover")
	rf := gi.NewFrame(discover)
	configDiscover(rf, mf)

	history := tabs.NewTab("History")
	ef := gi.NewFrame(history)
	configHistory(ef)

	gi.DefaultTopAppBar = func(tb *gi.TopAppBar) {
		gi.DefaultTopAppBarStd(tb)
		gi.NewButton(tb).SetIcon(icons.Add).SetText("New meal").OnClick(func(e events.Event) {
			newMeal(tb, mf, &osusu.Meal{})
		})
		gi.NewButton(tb).SetIcon(icons.Sort).SetText("Sort").OnClick(func(e events.Event) {
			d := gi.NewDialog(tb).Title("Sort and filter").FullWindow(true)
			giv.NewStructView(d).SetStruct(curOptions)
			d.OnAccept(func(e events.Event) {
				configSearch(mf)
				configHistory(ef)
				configDiscover(mf, rf)
			}).Run()
		})
	}

	gi.NewWindow(sc).SetNewWindow(false).Run()

	curGroup = &osusu.Group{}
	err := osusu.DB.First(curGroup, curUser.GroupID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			groups(sc)
		} else {
			gi.ErrorDialog(sc, err).Run()
		}
	}
}

func newMeal(ctx gi.Widget, mf *gi.Frame, meal *osusu.Meal) {
	d := gi.NewDialog(ctx).Title("Create meal").FullWindow(true)
	giv.NewStructView(d).SetStruct(meal)
	d.OnAccept(func(e events.Event) {
		err := osusu.DB.Create(meal).Error
		if err != nil {
			gi.ErrorDialog(d, err).Run()
			return
		}
		configSearch(mf)
	}).Cancel().Ok("Create").Run()
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
		switch w := w.(type) {
		case *gi.Label:
			w.Style(func(s *styles.Style) {
				s.SetNonSelectable()
				s.Align.Set(styles.AlignCenter)
				s.Text.Align = styles.AlignCenter
			})
		case *gi.Image:
			w.Style(func(s *styles.Style) {
				s.Max.Set(units.Em(20))
			})
		}
	})
}

func scoreGrid(card *gi.Frame, score *osusu.Score, showRecency bool) *gi.Layout {
	grid := gi.NewLayout(card)
	grid.Style(func(s *styles.Style) {
		s.Display = styles.DisplayGrid
		if showRecency {
			s.Columns = 6
		} else {
			s.Columns = 5
		}
		s.Align.X = styles.AlignCenter
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
	if grr.Log0(err) != nil {
		return nil
	}
	defer resp.Body.Close()
	img, _, err := images.Read(resp.Body)
	if grr.Log0(err) != nil {
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
