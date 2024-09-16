package main

import (
	"image"
	"net/http"
	"strconv"

	"cogentcore.org/core/base/errors"
	"cogentcore.org/core/base/iox/imagex"
	"cogentcore.org/core/colors"
	"cogentcore.org/core/core"
	"cogentcore.org/core/cursors"
	"cogentcore.org/core/enums"
	"cogentcore.org/core/events"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/abilities"
	"cogentcore.org/core/styles/units"
	"github.com/kkoreilly/osusu/osusu"
	"gorm.io/gorm"
)

var (
	curUser    *osusu.User
	curGroup   *osusu.Group
	curOptions = osusu.DefaultOptions()
)

func home() {
	b := core.NewBody("home")

	tabs := core.NewTabs(b).SetDeleteTabButtons(false)

	search := tabs.NewTab("Search")
	mf := core.NewFrame(search)
	configSearch(mf)

	discover := tabs.NewTab("Discover")
	rf := core.NewFrame(discover)
	tabs.Tabs().ChildByName("discover").(core.Widget).OnClick(func(e events.Event) {
		if !rf.HasChildren() {
			configDiscover(rf, mf)
		}
	})

	history := tabs.NewTab("History")
	ef := core.NewFrame(history)
	configHistory(ef)

	b.AddTopBar(func(pw core.Widget) {
		tb := b.DefaultTopAppBar(pw)
		core.NewButton(tb).SetIcon(icons.Add).SetText("New meal").OnClick(func(e events.Event) {
			newMeal(tb, mf, &osusu.Meal{})
		})
		core.NewButton(tb).SetIcon(icons.Sort).SetText("Sort").OnClick(func(e events.Event) {
			d := core.NewBody().AddTitle("Sort and filter")
			core.NewForm(d).SetStruct(curOptions)
			d.OnClose(func(e events.Event) {
				configSearch(mf)
				configHistory(ef)
				configDiscover(rf, mf)
			})
			d.NewFullDialog(tb).Run()
		})
	})

	b.NewWindow().SetNewWindow(false).Run()

	curGroup = &osusu.Group{}
	err := osusu.DB.First(curGroup, curUser.GroupID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			groups(b)
		} else {
			core.ErrorDialog(b, err)
		}
	}
}

func newMeal(ctx core.Widget, mf *core.Frame, meal *osusu.Meal) {
	d := core.NewBody().AddTitle("Create meal")
	core.NewForm(d).SetStruct(meal)
	d.AddBottomBar(func(pw core.Widget) {
		d.AddCancel(pw)
		d.AddOk(pw).SetText("Create").OnClick(func(e events.Event) {
			err := osusu.DB.Create(meal).Error
			if err != nil {
				core.ErrorDialog(d, err)
				return
			}
			configSearch(mf)
		})
	})
	d.NewFullDialog(ctx).Run()
}

func cardStyles(card *core.Frame) {
	card.Styler(func(s *styles.Style) {
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
	card.OnWidgetAdded(func(w core.Widget) {
		switch w := w.(type) {
		case *core.Text:
			w.Styler(func(s *styles.Style) {
				s.SetNonSelectable()
				s.Text.Align = styles.Center
			})
		case *core.Image:
			w.Styler(func(s *styles.Style) {
				s.Min.Set(units.Em(20))
				s.ObjectFit = styles.FitCover
			})
		}
	})
}

func scoreGrid(card *core.Frame, score *osusu.Score, showRecency bool) *core.Frame {
	grid := core.NewFrame(card)
	grid.Styler(func(s *styles.Style) {
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
		core.NewText(grid).SetType(core.TextLabelLarge).SetText(text)
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
		core.NewText(grid).SetText(strconv.Itoa(value))
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
	if errors.Log(err) != nil {
		return nil
	}
	defer resp.Body.Close()
	img, _, err := imagex.Read(resp.Body)
	if errors.Log(err) != nil {
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
