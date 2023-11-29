package main

import (
	"github.com/kkoreilly/osusu/osusu"
	"goki.dev/gi/v2/gi"
	"goki.dev/gi/v2/giv"
	"goki.dev/goosi/events"
)

func groups(b *gi.Body) {
	d := gi.NewBody("groups")
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
			gi.ErrorDialog(d, err)
			return
		}
		curGroup = newGroup
		curUser.GroupID = newGroup.ID
		err = osusu.DB.Save(curUser).Error
		if err != nil {
			gi.ErrorDialog(d, err)
		}
		d.Close()
	})
	d.NewFullDialog(b).Run()
}
