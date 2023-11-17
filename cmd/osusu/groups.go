package main

import (
	"github.com/kkoreilly/osusu/osusu"
	"goki.dev/gi/v2/gi"
	"goki.dev/gi/v2/giv"
	"goki.dev/goosi/events"
)

func groups(bsc *gi.Scene) {
	d := gi.NewDialog(bsc, "groups")
	gi.NewLabel(d).SetType(gi.LabelHeadlineMedium).SetText("Join group")
	groupCode := ""
	giv.NewValue(d, &groupCode)
	gi.NewButton(d).SetText("Join group")
	gi.NewLabel(d).SetType(gi.LabelHeadlineMedium).SetText("Create group")
	newGroup := &osusu.Group{OwnerID: curUser.ID, Owner: *curUser, Members: []osusu.User{*curUser}}
	giv.NewStructView(d).SetStruct(newGroup)
	gi.NewButton(d).SetText("Create group").OnClick(func(e events.Event) {
		err := osusu.DB.Create(newGroup).Error
		if err != nil {
			gi.NewDialog(d).Title("Error creating new group").Prompt(err.Error()).Run()
			return
		}
		curGroup = newGroup
		d.AcceptDialog()
	})
	d.Run()
}
