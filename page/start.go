package page

import (
	"github.com/kkoreilly/osusu/compo"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Start struct {
	app.Compo
}

func (s *Start) Render() app.UI {
	return &compo.Page{
		ID:                     "start",
		Title:                  "Start",
		Description:            "Osusu is an app for getting recommendations on what meals to eat in a group based on the ratings of each member of the group.",
		AuthenticationRequired: false,
		TitleElement:           "Welcome to Osusu!",
		SubtitleElement:        "Osusu is an app for getting recommendations on what meals to eat in a group based on the ratings of each member of the group.",
		Elements: []app.UI{
			&compo.ButtonRow{ID: "start-page", Buttons: []app.UI{
				&compo.Button{ID: "start-page-sign-in", Class: "secondary", Icon: "login", Text: "Sign In", OnClick: compo.NavigateEvent("/signin")},
				&compo.Button{ID: "start-page-sign-up", Class: "primary", Icon: "app_registration", Text: "Sign Up", OnClick: compo.NavigateEvent("/signup")},
			}},
			app.Div().ID("start-page-info").Body(
				StartPageInfos([]startPageInfo{
					{id: "recommendations", title: "Get Recommendations", body: "Get a ranked list of the best meals to eat based on your preferences. Each meal is scored on various metrics, so you can always pick the best meals that satisfy your needs.", img: "/web/images/recommendations.png"},
					{id: "everyone", title: "Works for Everyone", body: "Find meals that satisfy everyone's preferences and constraints, no matter how picky people are. Each person gets input on how they feel about each meal, and we find the meals that make everyone happy."},
					{id: "options", title: "Customize Options", body: "Customize your options and get recommendations that satisfy your needs. Whether you need something that is easy, cheap, and American for two people or something that is new, healthy, and Japanese for five people, you can get recommendations that fit whatever you need in the moment."},
					{id: "history", title: "Track History", body: "Track when you eat meals and how their quality changes over time. You can see how different attributes of a meal like taste, health, cost, and effort have gotten worse or better over time. Also, you can see when you ate each meal, how and for what meal you got it, and who you ate it with."},
				}),
			),
		},
	}
}

type startPageInfo struct {
	id    string
	title string
	body  string
	img   string
}

func StartPageInfos(infos []startPageInfo) app.UI {
	return app.Div().ID("start-page-infos").Body(
		app.Range(infos).Slice(func(i int) app.UI {
			info := infos[i]
			return app.Div().ID("start-page-info-container-"+info.id).Class("start-page-info-container").Body(
				app.Div().ID("start-page-info-text-container-"+info.id).Class("start-page-info-text-container").Body(
					app.H2().ID("start-page-info-title-"+info.id).Class("start-page-info-title").Text(info.title),
					app.P().ID("start-page-info-body-"+info.id).Class("start-page-info-body").Text(info.body),
				),
				app.Img().ID("start-page-info-image-"+info.id).Class("start-page-info-image").Src(info.img).Alt(info.title+" image"),
			)
		}),
	)
}
